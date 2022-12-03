## Configure aws CLI profile

    aws configure set aws_access_key_id "localstack" --profile localstack
    aws configure set aws_secret_access_key "localstack" --profile localstack
    aws configure set region "eu-central-1" --profile localstack
    aws configure set output "table" --profile test-profile

## Publish SQS message

    aws sns publish --endpoint-url=http://localhost:4566 --topic-arn arn:aws:sns:us-west-2:000000000000:import-zipcode-address --message '{ "zipCode": "01001000" }' --profile localstack --region us-west-2

    aws sns publish --endpoint-url=http://localhost:4566 --topic-arn arn:aws:sns:us-west-2:000000000000:import-zipcode-address --message '{ "zipCode": "49030340" }' --profile localstack --region us-west-2