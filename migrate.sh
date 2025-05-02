#!/bin/bash

# Set the migrations directory
MIGRATIONS_DIR="./sql/schema"

# Get the connection string 
CONN_STRING="postgres://postgres:postgres@localhost:5432/gator"

# Run the goose command from the migrations directory
(cd "$MIGRATIONS_DIR" && goose postgres "$CONN_STRING" "$@")