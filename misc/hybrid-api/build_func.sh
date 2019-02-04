#!/bin/sh

func=$1

base=$(pwd)


echo "Building $func lambda function"
cd $func
GOOS=linux go build -o $func main.go
zip -v $func.zip $func
cd $base

