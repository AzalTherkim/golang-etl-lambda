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

Using https://github.com/aws/aws-sdk-go to read and write to s3 buckets.

## How to deploy

For authentication to aws pls look at: https://registry.terraform.io/providers/hashicorp/aws/latest/docs

For how to build and get the zip file pls have a look at: https://docs.aws.amazon.com/lambda/latest/dg/golang-package.html

Example that works on windows:

### Build the go files:

```
set GOOS=linux
go build -o main main.go
```

Create the zip
```
%USERPROFILE%\Go\bin\build-lambda-zip.exe -output main.zip main
```

### Terraform
copy the main.zip into the terraform folder

change into the Terraform folder:
```
terraform init
```

```
terraform apply
```

# See it work

Upload a xml file to the xml bucket

# Thoughts

Not really happy with the writing to disk. As far as i can see only a move to the V2-SDK when that's available.