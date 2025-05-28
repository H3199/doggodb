package api

import (
	"context"
	"fmt"

	"github.com/H3199/doggodb/internal/api/generated/db"
	"github.com/H3199/doggodb/internal/data"
	"github.com/H3199/doggodb/internal/query"
)

type Server struct {
	db.UnimplementedDatabaseServiceServer
	Executor *query.Executor
}

func NewServer(executor *query.Executor) *Server {
	return &Server{
		Executor: executor,
	}
}

func (s *Server) Insert(ctx context.Context, req *db.InsertRequest) (*db.InsertResponse, error) {
	// Extract the table name
	tableName := req.GetTableName()

	// Prepare columns and values slices from the map
	var columns []string
	var values []string
	for col, val := range req.GetValues() {
		columns = append(columns, col)
		values = append(values, val)
	}

	// Construct the InsertStatement
	insertStmt := &query.InsertStatement{
		Table:   tableName,
		Columns: columns,
		Values:  values,
	}

	// Execute the InsertStatement
	_, err := s.Executor.Execute(insertStmt)
	if err != nil {
		return &db.InsertResponse{
			Success: false,
			Message: fmt.Sprintf("Insert failed: %v", err),
		}, nil
	}

	return &db.InsertResponse{
		Success: true,
		Message: "Insert successful",
	}, nil
}

func (s *Server) CreateTable(ctx context.Context, req *db.CreateTableRequest) (*db.CreateTableResponse, error) {
	stmt := &query.CreateTableStatement{
		Table: req.GetTableName(),
	}
	_, err := s.Executor.Execute(stmt)
	if err != nil {
		return &db.CreateTableResponse{
			Success: false,
			Message: err.Error(),
		}, nil
	}
	return &db.CreateTableResponse{
		Success: true,
		Message: "Table " + req.GetTableName() + " created successfully",
	}, nil
}

func (s *Server) Select(ctx context.Context, req *db.SelectRequest) (*db.SelectResponse, error) {
	selectStmt := &query.SelectStatement{
		Table:      req.GetTableName(),
		Columns:    req.GetColumns(),
		Conditions: req.GetConditions(),
	}

	rows, err := s.Executor.Execute(selectStmt)
	if err != nil {
		return nil, fmt.Errorf("failed to select: %w", err)
	}

	responseRows := []*db.Row{}
	for _, row := range rows.([]*data.Row) {
		values := make(map[string]string)
		for k, v := range row.Columns {
			// Assuming all values are string; convert if necessary.
			strVal, ok := v.(string)
			if !ok {
				strVal = fmt.Sprintf("%v", v) // fallback conversion
			}
			values[k] = strVal
		}
		responseRow := &db.Row{Values: values}
		responseRows = append(responseRows, responseRow)
	}

	return &db.SelectResponse{
		Rows: responseRows,
	}, nil
}

func (s *Server) Update(ctx context.Context, req *db.UpdateRequest) (*db.UpdateResponse, error) {
	updateStmt := &query.UpdateStatement{
		Table:       req.GetTableName(),
		Assignments: req.GetAssignments(),
		Conditions:  req.GetConditions(),
	}

	rowsUpdated, err := s.Executor.Execute(updateStmt)
	if err != nil {
		return nil, fmt.Errorf("failed to update: %w", err)
	}

	return &db.UpdateResponse{
		RowsUpdated: int32(rowsUpdated.(int)),
	}, nil
}

func (s *Server) Delete(ctx context.Context, req *db.DeleteRequest) (*db.DeleteResponse, error) {
	deleteStmt := &query.DeleteStatement{
		Table:      req.GetTableName(),
		Conditions: req.GetConditions(),
	}

	rowsDeleted, err := s.Executor.Execute(deleteStmt)
	if err != nil {
		return nil, fmt.Errorf("failed to delete: %w", err)
	}

	return &db.DeleteResponse{
		RowsDeleted: int32(rowsDeleted.(int)),
	}, nil
}
