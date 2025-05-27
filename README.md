# doggodb
A Go-based relational DB so simple, even a dog could've developed it!

![Alt text](doggo.png)

# Table of Contents

1. [Introduction](#introduction)
2. [Features](#features)
3. [Getting Started](#getting-started)
4. [API](#api)
5. [Design](#design)
6. [Testing](#testing)
   - [Unit Tests](#unit-tests)
   - [Integration Tests](#integration-tests)

## Introduction

Welcome to **Doggodb**, a relational database project written in Go. This project serves as both a personal learning journey and an exploration of deeper-level database mechanics, with the goal of improving proficiency in Go.

Currently in its early stages, Doggodb focuses on simplicity and supports basic CRUD operations. While minimalistic in functionality for now, the project lays the groundwork for more advanced features to be implemented in the future.

This project is not intended for production use but as an educational tool.

## Features

Doggodb is a lightweight relational database with the following features:

- **CRUD Operations**: Supports Create, Read, Update, and Delete functionalities.
- **gRPC API**: Exposes an interface for database operations using gRPC, making it accessible for distributed systems.
- **Custom Query Parser**: Implements a basic SQL-like query parser to handle database commands.
- **In-Memory Storage**: Uses in-memory storage for fast prototyping and testing.
- **Condition-Based Queries**: Allows filtering rows using `WHERE` conditions.
- **Unit & Integration Testing**: Includes extensive tests to ensure functionality and stability.

Planned features include:

- **TLS, IAM**: Some basic security features.
- **Transaction Support**: Add a transaction log to ensure atomicity and durability.
- **Command-Line Interface (CLI)**: Provide a user-friendly way to interact with the database directly through a terminal.
- **Enhanced SQL Syntax**: Support for more complex queries like joins, aggregations, and indexing.
- **Persistent Storage**: Extend the database to save data to disk for long-term storage.
- **Fancy features**: We could experiment with e.g. Bloom filters.
- **Clustering**: How hard could this be when we have a functioning JSON gRPC API? 

Doggodb aims to grow feature by feature, enhancing both functionality and understanding of database internals.

## Getting Started
Documentation under construction.

## API

### Overview

The gRPC API provides the core functionality of DoggoDB, enabling clients to perform CRUD operations on relational data. Each API endpoint corresponds to a specific database operation.

### Endpoints

#### 1. **CreateTable**
- **Request**:
```json
    {
    "table_name": "string"
    }
```
    
- **Response**:
```json
    {
    "success": "bool",
    "message": "string"
    }
```

- **Example**:
```bash
grpcurl -plaintext -d '{"table_name": "users"}' localhost:50051 doggodb.DatabaseService.CreateTable
```

#### 2. **Insert**:
- **Request**:
```json
    {
    "table_name": "string",
    "values": {
        "column1": "value1",
        "column2": "value2"
        }
    }
```

- **Response**:
```json
    {
    "success": "bool",
    "message": "string"
    }
```

- **Example**:
```bash
grpcurl -plaintext -d '{"table_name": "users", "values": {"id": "1", "name": "Alice"}}' localhost:50051 doggodb.DatabaseService.Insert
```

#### 3. **Select**:
- **Request**:
```json
    {
    "table_name": "string",
    "columns": ["column1", "column2"],
    "conditions": "string"
    }

```

- **Response**:
```json
{
  "rows": [
    {
      "values": {
        "column1": "value1",
        "column2": "value2"
      }
    },
  ]
}
```
- **Example**:
```bash
grpcurl -plaintext -d '{"table_name": "users", "columns": ["name"], "conditions": "id = 1"}' localhost:50051 doggodb.DatabaseService.Select
```

#### 4. **Update**:
- **Request**:
```json
    {
    "table_name": "string",
    "assignments": {
        "column1": "new_value1",
        "column2": "new_value2"
        },
    "conditions": "string"
    }
```

- **Response**:
```json
{
  "rows_updated": "int"
}
```

- **Example**:
```bash
grpcurl -plaintext -d '{"table_name": "users", "assignments": {"name": "Alice Updated"}, "conditions": "id = 1"}' localhost:50051 doggodb.DatabaseService.Update

```

#### 5. **Delete**:
- **Request**:
```json
    {
    "table_name": "string",
    "conditions": "string"
    }
```

- **Response**:
```json
{
  "rows_deleted": "int"
}
```

- **Example**:
```bash
grpcurl -plaintext -d '{"table_name": "users", "conditions": "id = 1"}' localhost:50051 doggodb.DatabaseService.Delete

```

## Design
Documentation under construction.

## Testing

- **Unit testing**:
```bash
go test test/unit/* -v

```

- **Integration testing**:
```bash
go test test/integration/* -v

```





