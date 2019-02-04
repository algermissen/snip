hybrid-api
==========

A proof of concept for a hybrid API of REST and other integration styles.

Preparations in AWS
===================

- Create role for lambda execution

    arn:aws:iam::$account:role/lambda_basic_execution


Building the Tools and Lambda Functions
=======================================

Make sure your GOPATH is set as you want and then run

    go get github.com/aws/aws-lambda-go/lambda
    go get github.com/aws/aws-lambda-go/events
    go get github.com/aws/aws-sdk-go/aws/session
    go get github.com/aws/aws-sdk-go/aws
    go get github.com/aws/aws-sdk-go/service/kinesis
    go get github.com/google/uuid


Next, run

    ./build.sh


Create/Update Lambda Functions
==============================

Function <func> must be in directory called <func>.

./create_lambda <func>
