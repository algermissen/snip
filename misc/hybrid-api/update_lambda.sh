#!/bin/sh

account=627211582084

func=$1

aws lambda update-function-code \
  --function-name $func \
  --region eu-central-1 \
  --zip-file fileb:///Users/janalgermissen/projects/hse/dataload/hybrid-api/$func/$func.zip 
