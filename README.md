# golang-etl-lambda
short poc for an s3 etl with lambda

# Setup
s3 bucket with xml data -> lambda function that transforms to json -> s3 bucket with json output

# Infrastructure with Terraform

Create 2 S3 Buckets

Create Lambda Function
* needs access to S3
* and lambda execution role

# Go for the lambda function

Using https://github.com/clbanning/mxj to transform from xml to json.
Using https://github.com/aws/aws-sdk-go to read and write to s3 buckets