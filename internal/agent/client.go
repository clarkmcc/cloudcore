// MIT License
//
// Copyright (c) 2024 Clark McCauley
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

package agent

import (
	"context"
	"crypto/tls"
	"errors"
	"github.com/clarkmcc/brpc"
	"github.com/clarkmcc/cloudcore/internal/rpc"
	"github.com/spf13/cast"
	"github.com/spf13/cobra"
	"go.uber.org/multierr"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/connectivity"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"gopkg.in/tomb.v2"
	"io"
	"sync"
	"time"
)

type Client struct {
	dialer   func(ctx context.Context) (*brpc.ClientConn, error)
	service  rpc.AgentServer
	logger   *zap.Logger
	tomb     *tomb.Tomb
	shutdown <-chan struct{}

	tokenManager *tokenManager

	lock      sync.Mutex
	conn      *brpc.ClientConn
	auth      rpc.AuthenticationClient
	agent     rpc.AgentManagerClient
	notify    rpc.AgentManager_NotificationsClient
	resetConn bool
}

func (c *Client) Ping(ctx context.Context) error {
	c.lock.Lock()
	defer c.lock.Unlock()

	client, err := c.getAuthClientLocked(ctx)
	if err != nil {
		return err
	}
	_, err = client.Ping(ctx, &rpc.PingRequest{})
	return err
}

func (c *Client) Notify(ctx context.Context, notification *rpc.ClientNotification) error {
	c.lock.Lock()
	defer c.lock.Unlock()

	ctx, err := c.getAuthContextLocked(ctx)
	if err != nil {
		return err
	}
	stream, err := c.getNotificationsStreamLocked(ctx)
	if err != nil {
		return err
	}
	err = stream.Send(notification)
	if err != nil {
		if errors.Is(err, io.EOF) {
			c.notify = nil
		}
		return err
	}
	return nil
}

func (c *Client) getNotificationsStreamLocked(ctx context.Context) (rpc.AgentManager_NotificationsClient, error) {
	err := c.connectStreamsLocked(ctx)
	if err != nil {
		return nil, err
	}
	return c.notify, nil
}

func (c *Client) getAgentClientLocked(ctx context.Context) (rpc.AgentManagerClient, error) {
	err := c.connectLocked(ctx)
	if err != nil {
		return nil, err
	}
	return c.agent, nil
}

func (c *Client) getAuthContextLocked(ctx context.Context) (context.Context, error) {
	tk, err := c.tokenManager.getAuthTokenLocked(ctx)
	if err != nil {
		if s, ok := status.FromError(err); ok {
			if s.Code() == codes.Unavailable {
				c.resetConn = true
			}
		}
		return nil, err
	}
	return metadata.NewOutgoingContext(ctx, metadata.Pairs("token", tk)), nil
}

// getAuthClientLocked returns a connected authentication client. This function should only be
// called from the token manager, or from the Ping method above. It should only be called
// in contexts where the client is already locked.
//
// Specifically, when this function is called from the token manager, the token manager should
// already be locked because the token manager's getAuthTokenLocked method is only ever called from
// one of the client's public RPC methods which actually acquire the lock.
func (c *Client) getAuthClientLocked(ctx context.Context) (rpc.AuthenticationClient, error) {
	err := c.connectLocked(ctx)
	if err != nil {
		return nil, err
	}
	return c.auth, nil
}

// connectLocked connects to the server if not already connected. If the client is shutdown,
// then this function will attempt a reconnect.
func (c *Client) connectLocked(ctx context.Context) error {
	if c.conn != nil && c.conn.ClientConn != nil && !c.resetConn {
		switch v := c.conn.GetState(); v {
		case connectivity.Shutdown:
			fallthrough
		case connectivity.TransientFailure:
			return c.setupClientsLocked(ctx)
		default:
			_ = v
			return nil
		}
	}
	return c.setupClientsLocked(ctx)
}

func (c *Client) connectStreamsLocked(ctx context.Context) (err error) {
	ctx, err = c.getAuthContextLocked(ctx)
	if err != nil {
		return err
	}

	if c.notify != nil {
		err := c.notify.Context().Err()
		if err == nil {
			return nil
		}
	}
	c.notify, err = c.agent.Notifications(c.shutdownCtx(ctx))
	if err != nil {
		return err
	}
	return nil
}

func (c *Client) setupClientsLocked(ctx context.Context) (err error) {
	ctx = c.shutdownCtx(ctx)
	c.conn, err = c.dialer(ctx)
	if err != nil {
		return err
	}
	c.auth = rpc.NewAuthenticationClient(c.conn)
	c.agent = rpc.NewAgentManagerClient(c.conn)

	c.tomb.Go(func() error {
		err = brpc.ServeClientService[rpc.AgentServer](ctx.Done(), c.conn, func(registrar grpc.ServiceRegistrar) {
			rpc.RegisterAgentServer(registrar, c.service)
		})
		if err != nil {
			c.logger.Error("error serving client service", zap.Error(err))
		}
		return nil
	})

	// Reset the streams if we're reconnecting.
	if c.notify != nil {
		multierr.AppendFunc(&err, c.notify.CloseSend)
	}
	c.notify = nil
	return nil
}

// shutdownCtx returns a special context that cancels once the tomb is fully
// dead. The idea here is that the client should be the last thing to shut
// down because some shutdown procedures need a working client in order to
// work (such as sending shutdown notifications).
//
// If the provided context dies, we wait for 5 seconds for the tomb to fully
// die before cancelling the return context. The idea is we:
//  1. Don't want to block forever if the provided context is cancelled but
//     the tomb doesn't die.
//  2. Don't want to cancel the returned context at the same time we get the
//     shutdown signal.
func (c *Client) shutdownCtx(ctx context.Context) context.Context {
	ctx, cancel := context.WithCancel(context.WithoutCancel(ctx))
	go func() {
		select {
		case <-ctx.Done():
		case <-c.shutdown:
		}
		time.Sleep(5 * time.Second)
		c.logger.Error("shutting down client context")
		cancel()
	}()
	return ctx
}

func NewClient(
	ctx context.Context,
	config *Config,
	tomb *tomb.Tomb,
	cmd *cobra.Command,
	db Database,
	logger *zap.Logger,
	service rpc.AgentServer,
	metadataProvider SystemMetadataProvider,
) *Client {
	c := &Client{
		tomb:     tomb,
		service:  service,
		shutdown: ctx.Done(),
		logger:   logger.Named("client"),
		dialer: func(ctx context.Context) (*brpc.ClientConn, error) {
			return brpc.DialContext(ctx, config.Server.Endpoint, &tls.Config{
				InsecureSkipVerify: cast.ToBool(cmd.Flag("insecure-skip-verify").Value.String()),
			})
		},
	}
	c.tokenManager = newTokenManager(config, db, logger, c, metadataProvider)
	return c
}
