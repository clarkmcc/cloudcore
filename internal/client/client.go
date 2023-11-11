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
	"github.com/clarkmcc/cloudcore/internal/agentdb"
	"github.com/clarkmcc/cloudcore/internal/config"
	"github.com/clarkmcc/cloudcore/internal/rpc"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"sync"
)

type Client struct {
	dialer func(ctx context.Context) (*grpc.ClientConn, error)

	tokenManager *tokenManager

	lock  sync.Mutex
	conn  *grpc.ClientConn
	auth  rpc.AuthenticationClient
	agent rpc.AgentManagerClient
}

func (c *Client) Ping(ctx context.Context) error {
	client, err := c.getAuthClient(ctx)
	if err != nil {
		return err
	}
	_, err = client.Ping(ctx, &rpc.PingRequest{})
	return err
}

func (c *Client) UploadMetadata(ctx context.Context, metadata *rpc.SystemMetadata) error {
	ctx, err := c.getAuthContext(ctx)
	if err != nil {
		return err
	}
	client, err := c.getAgentClient(ctx)
	if err != nil {
		return err
	}
	_, err = client.UploadMetadata(ctx, &rpc.UploadMetadataRequest{
		SystemMetadata: metadata,
	})
	return err
}

func (c *Client) getAuthClient(ctx context.Context) (rpc.AuthenticationClient, error) {
	err := c.connect(ctx)
	if err != nil {
		return nil, err
	}
	return c.auth, nil
}

func (c *Client) getAgentClient(ctx context.Context) (rpc.AgentManagerClient, error) {
	err := c.connect(ctx)
	if err != nil {
		return nil, err
	}
	return c.agent, nil
}

func (c *Client) getAuthContext(ctx context.Context) (context.Context, error) {
	tk, err := c.tokenManager.getAuthToken(ctx)
	if err != nil {
		return nil, err
	}
	return metadata.NewOutgoingContext(ctx, metadata.Pairs("token", tk)), nil
}

// connect connects to the server if not already connected. If the client is shutdown,
// then this function will attempt a reconnect.
func (c *Client) connect(ctx context.Context) error {
	c.lock.Lock()
	defer c.lock.Unlock()
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

func (c *Client) setupClientsLocked(ctx context.Context) (err error) {
	c.conn, err = c.dialer(ctx)
	if err != nil {
		return err
	}
	c.auth = rpc.NewAuthenticationClient(c.conn)
	c.agent = rpc.NewAgentManagerClient(c.conn)
	return nil
}

func New(config *config.AgentConfig, db agentdb.AgentDB, logger *zap.Logger) *Client {
	c := &Client{
		dialer: func(ctx context.Context) (*grpc.ClientConn, error) {
			return grpc.Dial(config.Server.Endpoint, grpc.WithTransportCredentials(insecure.NewCredentials()))
		},
	}
	c.tokenManager = newTokenManager(config, db, logger, c)
	return c
}
