#!/bin/bash

export AWS_ACCESS_KEY_ID=localstack
export AWS_SECRET_ACCESS_KEY=localstack

echo "installing jq"
apk update && apk add --no-cache jq

echo "configure region [us-west-2]"
aws configure set default.region us-west-2

echo "configure sns topic"
TOPIC_NAME="import-zipcode-address"
TOPIC_ARN=$(aws --endpoint-url http://localhost:4566 sns create-topic --output text --name "$TOPIC_NAME")

echo "configure sqs queue"
QUEUE_NAME="import-zipcode-address"
QUEUE_URL=$(aws --endpoint-url http://localhost:4566 sqs create-queue --queue-name "$QUEUE_NAME" --output text)
QUEUE_ARN=$(aws --endpoint-url http://localhost:4566 sqs get-queue-attributes --queue-url "$QUEUE_URL" | jq -r ".Attributes.QueueArn")

echo "configure sqs subscription"
aws --endpoint-url http://localhost:4566 sns subscribe --topic-arn "$TOPIC_ARN" --protocol sqs --notification-endpoint "$QUEUE_ARN" --output text

aws --endpoint-url http://localhost:4566 sns list-subscriptions
curl http://localhost:4566/health