package domain

import (
	"context"
	"errors"
	"io"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Client struct {
	conn          *grpc.ClientConn
	productClient ProductServiceClient
}

func NewClient(addr string) (*Client, error) {
	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	return &Client{
		conn:          conn,
		productClient: NewProductServiceClient(conn),
	}, nil
}

func (c *Client) CloseConnection() error {
	return c.conn.Close()
}

func (c *Client) Fetch(ctx context.Context, req *FetchRequest) error {
	_, err := c.productClient.Fetch(ctx, req)

	return err
}

func (c *Client) List(ctx context.Context, req []*Filters, res chan<- *ListResponse) error {
	stream, err := c.productClient.List(ctx)
	if err != nil {
		return err
	}

	waitc := make(chan struct{})

	go func() {
		for {
			in, err := stream.Recv()
			if errors.Is(err, io.EOF) {
				// read done.
				close(waitc)
				close(res)

				return
			}

			if err != nil {
				log.Fatalf("Failed to receive a note : %v", err)
			}

			res <- in
		}
	}()

	for _, filter := range req {
		if err = stream.Send(filter); err != nil {
			return err
		}
	}

	if err = stream.CloseSend(); err != nil {
		return err
	}

	<-waitc

	return nil
}
