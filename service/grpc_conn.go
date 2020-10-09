//    Copyright 2017 Ewout Prangsma
//
//    Licensed under the Apache License, Version 2.0 (the "License");
//    you may not use this file except in compliance with the License.
//    You may obtain a copy of the License at
//
//        http://www.apache.org/licenses/LICENSE-2.0
//
//    Unless required by applicable law or agreed to in writing, software
//    distributed under the License is distributed on an "AS IS" BASIS,
//    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//    See the License for the specific language governing permissions and
//    limitations under the License.

package service

import (
	"context"

	"google.golang.org/grpc"
)

type grpcConn struct {
	conn   *grpc.ClientConn
	ctx    context.Context
	cancel context.CancelFunc
}

func (c *grpcConn) Close() {
	if c.cancel != nil {
		c.cancel()
		c.cancel = nil
	}
	if c.conn != nil {
		c.conn.Close()
		c.conn = nil
	}
	c.ctx = nil
}

func (c *grpcConn) SetConn(ctx context.Context, conn *grpc.ClientConn) {
	c.Close()
	c.ctx, c.cancel = context.WithCancel(ctx)
	c.conn = conn
}
