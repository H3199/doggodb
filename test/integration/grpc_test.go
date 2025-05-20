package server_test

import (
	"context"
	"net"
	"testing"

	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"

	"github.com/H3199/doggodb/internal/api"
	"github.com/H3199/doggodb/internal/api/generated/db"
	"github.com/H3199/doggodb/internal/data"
	"github.com/H3199/doggodb/internal/query"
)

func startTestGRPCServer(t *testing.T) (db.DatabaseServiceClient, func()) {
	t.Helper()

	storage := data.NewInMemoryStorage()
	executor := query.NewExecutor(*storage)

	grpcServer := grpc.NewServer()
	srv := &api.Server{Executor: executor}
	db.RegisterDatabaseServiceServer(grpcServer, srv)

	lis, err := net.Listen("tcp", "localhost:0")
	require.NoError(t, err)

	go grpcServer.Serve(lis)

	conn, err := grpc.Dial(lis.Addr().String(), grpc.WithInsecure())
	require.NoError(t, err)

	client := db.NewDatabaseServiceClient(conn)

	cleanup := func() {
		grpcServer.Stop()
		conn.Close()
		lis.Close()
	}

	return client, cleanup
}

func TestGRPCFullFlow(t *testing.T) {
	client, cleanup := startTestGRPCServer(t)
	defer cleanup()

	ctx := context.Background()

	// 1. Create Table
	t.Log("Creating table 'users'")
	createResp, err := client.CreateTable(ctx, &db.CreateTableRequest{TableName: "users"})
	require.NoError(t, err)
	require.True(t, createResp.Success, createResp.Message)

	// 2. Insert a row
	t.Log("Inserting a user row")
	insertResp, err := client.Insert(ctx, &db.InsertRequest{
		TableName: "users",
		Values: map[string]string{
			"id":    "1",
			"name":  "Alice",
			"email": "alice@example.com",
		},
	})
	require.NoError(t, err)
	require.True(t, insertResp.Success, insertResp.Message)

	// 3. Select the inserted row
	t.Log("Selecting all columns from 'users'")
	selectResp, err := client.Select(ctx, &db.SelectRequest{
		TableName:  "users",
		Columns:    []string{"*"}, // empty slice = SELECT *
		Conditions: "",
	})
	require.NoError(t, err)
	require.Len(t, selectResp.Rows, 1)
	require.Equal(t, "Alice", selectResp.Rows[0].Values["name"])
	require.Equal(t, "alice@example.com", selectResp.Rows[0].Values["email"])

	// 4. Update the row
	t.Log("Updating user row")
	updateResp, err := client.Update(ctx, &db.UpdateRequest{
		TableName:   "users",
		Assignments: map[string]string{"email": "alice@newdomain.com"},
		Conditions:  "id = '1'",
	})
	require.NoError(t, err)
	require.EqualValues(t, 1, updateResp.RowsUpdated)

	// 5. Select again and check update
	t.Log("Selecting updated row")
	selectResp2, err := client.Select(ctx, &db.SelectRequest{
		TableName:  "users",
		Columns:    []string{"id", "name", "email"},
		Conditions: "id = '1'",
	})
	require.NoError(t, err)
	require.Len(t, selectResp2.Rows, 1)
	require.Equal(t, "alice@newdomain.com", selectResp2.Rows[0].Values["email"])

	// 6. Delete the row
	t.Log("Deleting user row")
	deleteResp, err := client.Delete(ctx, &db.DeleteRequest{
		TableName:  "users",
		Conditions: "id = '1'",
	})
	require.NoError(t, err)
	require.EqualValues(t, 1, deleteResp.RowsDeleted)

	// 7. Select again to confirm deletion
	t.Log("Selecting after deletion to confirm no rows")
	selectResp3, err := client.Select(ctx, &db.SelectRequest{
		TableName:  "users",
		Columns:    []string{},
		Conditions: "",
	})
	require.NoError(t, err)
	require.Len(t, selectResp3.Rows, 0)
}
