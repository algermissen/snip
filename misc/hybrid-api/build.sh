#!/bin/sh

base=$(pwd)

echo "Building producer"
cd producer
go build -o producer main.go
cd $base

echo "Building kinesis-consumer"
cd kinesis-consumer
go build -o kinesis-consumer main.go
cd $base

./build_func.sh home

./build_func.sh index

./build_func.sh wsdl

./build_func.sh soap

./build_func.sh http

