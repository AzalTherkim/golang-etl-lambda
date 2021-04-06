terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 3.0"
    }
  }
}

# Configure the AWS Provider
provider "aws" {
  region = "eu-central-1"
}

resource "aws_iam_role" "iam_for_lambda" {
  name = "iam_for_lambda"

  assume_role_policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Action": "sts:AssumeRole",
      "Principal": {
        "Service": "lambda.amazonaws.com"
      },
      "Effect": "Allow"
    }
  ]
}
EOF
}

resource "aws_iam_role_policy_attachment" "lambda_execution" {

  role       = aws_iam_role.iam_for_lambda.name
  policy_arn = "arn:aws:iam::aws:policy/service-role/AWSLambdaBasicExecutionRole"

}

# This is very broad might make more sense to change this in the long term
resource "aws_iam_role_policy_attachment" "s3_access" {

  role       = aws_iam_role.iam_for_lambda.name
  policy_arn = "arn:aws:iam::aws:policy/AmazonS3FullAccess"

}


resource "aws_lambda_permission" "allow_xml_bucket" {
  statement_id  = "AllowExecutionFromS3Bucket"
  action        = "lambda:InvokeFunction"
  function_name = aws_lambda_function.transform_code.arn
  principal     = "s3.amazonaws.com"
  source_arn    = aws_s3_bucket.xml_data.arn
}

resource "aws_lambda_function" "transform_code" {
  filename      = "main.zip"
  function_name = var.function_name
  role          = aws_iam_role.iam_for_lambda.arn
  handler       = "main"
  runtime       = "go1.x" 
  environment {
    variables = {
      OUTPUT_BUCKET = aws_s3_bucket.json_data.id
    }
  }

}



data "aws_canonical_user_id" "current_user" {}

resource "aws_s3_bucket" "xml_data" {
  bucket = var.xml_bucket

  grant {
    id          = data.aws_canonical_user_id.current_user.id
    type        = "CanonicalUser"
    permissions = ["FULL_CONTROL"]
  }
}

resource "aws_s3_bucket" "json_data" {
  bucket = var.json_bucket

  grant {
    id          = data.aws_canonical_user_id.current_user.id
    type        = "CanonicalUser"
    permissions = ["FULL_CONTROL"]
  }

}

resource "aws_s3_bucket_notification" "bucket_notification" {
  bucket = aws_s3_bucket.xml_data.id

  lambda_function {
    lambda_function_arn = aws_lambda_function.transform_code.arn
    events              = ["s3:ObjectCreated:*"]
  }

  depends_on = [aws_lambda_permission.allow_xml_bucket]
}