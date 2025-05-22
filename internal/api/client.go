package api

import (
	"context"
	"fmt"

	"github.com/H3199/doggodb/internal/api/generated/db"
	"google.golang.org/grpc"
)

type Client struct {
	conn   *grpc.ClientConn
	client db.DatabaseServiceClient
}

func NewClient(address string) (*Client, error) {
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		return nil, fmt.Errorf("failed to connect: %w", err)
	}
	return &Client{
		conn:   conn,
		client: db.NewDatabaseServiceClient(conn),
	}, nil
}

func (c *Client) CreateTable(tableName string) error {
	resp, err := c.client.CreateTable(context.Background(), &db.CreateTableRequest{
		TableName: tableName,
	})
	if err != nil {
		return fmt.Errorf("failed to create table: %w", err)
	}
	if !resp.Success {
		return fmt.Errorf("server error: %s", resp.Message)
	}
	return nil
}

func (c *Client) Insert(tableName string, values map[string]string) error {
	resp, err := c.client.Insert(context.Background(), &db.InsertRequest{
		TableName: tableName,
		Values:    values,
	})
	if err != nil {
		return fmt.Errorf("failed to insert data: %w", err)
	}
	if !resp.Success {
		return fmt.Errorf("server error: %s", resp.Message)
	}
	return nil
}

func (c *Client) Select(tableName string, columns []string, conditions string) ([]*db.Row, error) {
	resp, err := c.client.Select(context.Background(), &db.SelectRequest{
		TableName:  tableName,
		Columns:    columns,
		Conditions: conditions,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to select data: %w", err)
	}
	return resp.Rows, nil
}

func (c *Client) Close() {
	c.conn.Close()
}
