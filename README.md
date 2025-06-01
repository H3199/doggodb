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
- **Types**: Support for typed values
- **Dockerfile**
- **Kubernetes operator**: This should be interesting

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
## Storage Driver Architecture

### Overview
The doggodb storage driver is designed around two core concepts:

1. **B-Tree Structured Storage Using the Filesystem**  
2. **Write-Ahead Logging (WAL) for Durability and Recovery**

These components work together to provide efficient, durable, and crash-resilient storage for the database.
NOTE & TODO: none of this have been implemented yet. This is just the architechture design.

---

### 1. Filesystem-Based B-Tree Storage

#### Concept
- The storage driver uses the **native filesystem directory and file hierarchy** to represent the structure of a B-tree index for each table.
- Each **node in the B-tree** is represented by a directory (internal node) or a file (leaf node).
- **Keys** are stored as files containing serialized row data (`JSON` format).
- **Child nodes** are subdirectories representing branches of the B-tree.

#### Example Structure for a Table `users`
```
users/               # Root directory for the 'users' table
├── key_2.json       # Leaf node file with key=2 row data
├── child_1/         # Subdirectory representing a B-tree child node
│   ├── key_1.json
│   └── key_3.json
└── child_2/
    ├── key_4.json
    └── key_5.json
```

#### Benefits
- Navigable by shell and other tools, enabling manual inspection and lightweight scripting.
- Avoids reinventing low-level file management.
- Intuitive mapping between logical B-tree nodes and physical filesystem entities.

---

### 2. Write-Ahead Log (WAL)

#### Concept
- The WAL is an **append-only log** file that records every change operation (inserts, updates, deletes) before applying it to the B-tree filesystem.
- Acts as a durable, sequential record of all intended changes.

#### Functionality
- **Before applying changes** to the filesystem B-tree, the operation is logged to the WAL.
- WAL entries are structured as JSON-encoded instructions, e.g.:
```json
{"type": "INSERT", "key": 6, "row": {"ID": 6, "Name": "Fiona", "Age": 23}}
```
- In case of crashes or failures, the WAL can be **replayed** to restore consistency by reapplying all logged operations.

#### Checkpointing
- Periodically, changes from WAL are **flushed** to the filesystem B-tree (actual files and directories).
- After successful flush, the WAL file is truncated or rotated to avoid unbounded growth.

---

### 3. Decoupling of In-Memory and On-Disk Storage

#### In-Memory Storage
- Serves as the **fast, transient data structure** for query processing and data manipulation.
- All SQL commands operate first on this structure for responsiveness.

#### On-Disk Storage
- Represented by the **filesystem-based B-tree** directories and files.
- Holds the durable, persistent state of the database.

#### Synchronization Workflow
1. **SQL Command Execution**:
   - Modify the in-memory structure.
   - Append the operation to the WAL for durability.

2. **Checkpointing**:
   - Periodically write in-memory changes to the filesystem B-tree.
   - Truncate or clear the WAL once changes are safely persisted.

3. **Crash Recovery**:
   - On startup, replay the WAL to apply any unpersisted operations to the filesystem structure.
   - This restores on-disk state to the last consistent point.

---

### 4. Future Considerations

- **Concurrency and Locking**: Mechanisms to handle concurrent access and updates safely.
- **Node Rebalancing**: Implementing B-tree balancing in the filesystem hierarchy.
- **Compaction**: Cleaning up stale files and optimizing storage usage.
- **Performance**: Caching strategies and WAL optimization for faster commits.

## Testing

- **Unit testing**:
```bash
go test test/unit/* -v

```

- **Integration testing**:
```bash
go test test/integration/* -v

```





