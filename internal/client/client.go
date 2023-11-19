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

package client

import (
	"context"
	"errors"
	"github.com/clarkmcc/cloudcore/cmd/cloudcored/config"
	"github.com/clarkmcc/cloudcore/internal/agentdb"
	"github.com/clarkmcc/cloudcore/internal/rpc"
	"go.uber.org/multierr"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"io"
	"sync"
)

type Client struct {
	dialer func(ctx context.Context) (*grpc.ClientConn, error)

	tokenManager *tokenManager

	lock   sync.Mutex
	conn   *grpc.ClientConn
	auth   rpc.AuthenticationClient
	agent  rpc.AgentManagerClient
	notify rpc.AgentManager_NotificationsClient
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

func (c *Client) UploadMetadata(ctx context.Context, metadata *rpc.SystemMetadata) (*rpc.UploadMetadataResponse, error) {
	c.lock.Lock()
	defer c.lock.Unlock()

	ctx, err := c.getAuthContextLocked(ctx)
	if err != nil {
		return nil, err
	}
	client, err := c.getAgentClientLocked(ctx)
	if err != nil {
		return nil, err
	}
	return client.UploadMetadata(ctx, &rpc.UploadMetadataRequest{
		SystemMetadata: metadata,
	})
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
	if c.conn != nil {
		switch c.conn.GetState() {
		case connectivity.Shutdown:
			return c.setupClientsLocked(ctx)
		default:
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
	c.notify, err = c.agent.Notifications(ctx)
	if err != nil {
		return err
	}
	return nil
}

func (c *Client) setupClientsLocked(ctx context.Context) (err error) {
	c.conn, err = c.dialer(ctx)
	if err != nil {
		return err
	}
	c.auth = rpc.NewAuthenticationClient(c.conn)
	c.agent = rpc.NewAgentManagerClient(c.conn)

	// Reset the streams if we're reconnecting.
	if c.notify != nil {
		multierr.AppendFunc(&err, c.notify.CloseSend)
	}
	c.notify = nil
	return nil
}

func New(config *config.Config, db agentdb.AgentDB, logger *zap.Logger) *Client {
	c := &Client{
		dialer: func(ctx context.Context) (*grpc.ClientConn, error) {
			return grpc.Dial(config.Server.Endpoint, grpc.WithTransportCredentials(insecure.NewCredentials()))
		},
	}
	c.tokenManager = newTokenManager(config, db, logger, c)
	return c
}
