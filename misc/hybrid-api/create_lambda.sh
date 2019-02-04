#!/bin/sh

account=627211582084

func=$1


aws lambda create-function \
  --function-name $func \
  --memory 128 \
  --runtime go1.x \
  --region eu-central-1 \
  --role arn:aws:iam::$account:role/lambda_basic_execution \
  --zip-file fileb:///Users/janalgermissen/projects/hse/dataload/hybrid-api/$func/$func.zip \
  --handler $func
