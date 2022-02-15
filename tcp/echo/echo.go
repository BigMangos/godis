// Package echo for test
// A echo server to test whether the server is functioning normally
package echo

import (
	"bufio"
	"context"
	"io"
	"net"
	"sync"
	"time"

	"github.com/hdt3213/godis/lib/logger"
	"github.com/hdt3213/godis/lib/sync/atomic"
	"github.com/hdt3213/godis/lib/sync/wait"
)

// Handler echos received line to client, using for test
type Handler struct {
	activeConn sync.Map
	closing    atomic.Boolean
}

// NewEchoHandler creates Handler
func NewEchoHandler() *Handler {
	return &Handler{}
}

// Client is client for Handler, using for test
type Client struct {
	Conn    net.Conn
	Waiting wait.Wait
}

// Close connection
func (c *Client) Close() error {
	c.Waiting.WaitWithTimeout(10 * time.Second)
	c.Conn.Close()
	return nil
}

// Handle echos received line to client
func (h *Handler) Handle(ctx context.Context, conn net.Conn) {
	if h.closing.Get() {
		// closing handler refuse new connection
		_ = conn.Close()
	}

	client := &Client{
		Conn: conn,
	}
	h.activeConn.Store(client, struct{}{})

	reader := bufio.NewReader(conn)
	for {
		// may occurs: client EOF, client timeout, server early close
		msg, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				logger.Info("connection close")
				h.activeConn.Delete(client)
			} else {
				logger.Warn(err)
			}
			return
		}
		client.Waiting.Add(1)
		//logger.Info("sleeping")
		//time.Sleep(10 * time.Second)
		b := []byte(msg)
		_, _ = conn.Write(b)
		client.Waiting.Done()
	}
}

// Close stops echo handler
func (h *Handler) Close() error {
	logger.Info("handler shutting down...")
	h.closing.Set(true)
	h.activeConn.Range(func(key interface{}, val interface{}) bool {
		client := key.(*Client)
		_ = client.Close()
		return true
	})
	return nil
}
