#!/bin/bash

echo "Running swag init..."
swag init

if [ $? -ne 0 ]; then
    echo "swag init failed with error code $?."
    exit 1
fi

echo "swag init completed successfully."

echo "Running go run main.go..."
go run main.go

if [ $? -ne 0 ]; then
    echo "go run main.go failed with error code $?."
    exit 1
fi

