#!/bin/ksh

echo "--- Cleaning"
rm w-cli
rm w-svr

echo "--- Building Warehouse Server"
go build -o w-svr server/main/main.go

echo "--- Building CLI"
go build -o w-cli tui/main/main.go
