locals {
  name = "linker"
  target_file = "deploy.zip"
}

resource "aws_lambda_function" "lambda" {
  filename      = local.target_file
  function_name = local.name
  role          = aws_iam_role.iam_role.arn
  handler       = "main"
  runtime       = "go1.x"

  source_code_hash = filebase64sha256(local.target_file)

  environment {
    variables = {
      DYNAMODB_TABLE_NAME = aws_dynamodb_table.dynamodb.name
    }
  }

  depends_on = [
    aws_iam_role_policy_attachment.logging,
    aws_iam_role_policy_attachment.dynamodb,
  ]
}

resource "aws_iam_role" "iam_role" {
  name = "lambda_${local.name}"

  assume_role_policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Action": "sts:AssumeRole",
      "Principal": {
        "Service": "lambda.amazonaws.com"
      },
      "Effect": "Allow",
      "Sid": ""
    }
  ]
}
EOF
}

resource "aws_cloudwatch_log_group" "logging" {
  name              = "/aws/lambda/${local.name}"
  retention_in_days = 14
}

data "aws_iam_policy_document" "logging" {
  statement {
    actions   = [
        "logs:CreateLogGroup",
        "logs:CreateLogStream",
        "logs:PutLogEvents"
    ]
    resources = [aws_cloudwatch_log_group.logging.arn]
  }
}

resource "aws_iam_policy" "logging" {
  name        = "lambda_${local.name}_logging"
  path        = "/"
  description = "IAM policy to allow lambda to write logs on CloudWatch"

  policy = data.aws_iam_policy_document.logging.json
}

resource "aws_iam_role_policy_attachment" "logging" {
  role       = aws_iam_role.iam_role.name
  policy_arn = aws_iam_policy.logging.arn
}


data "aws_iam_policy_document" "dynamodb" {
  statement {
    actions   = ["dynamodb:*"]
    resources = [aws_dynamodb_table.dynamodb.arn]
  }
}

resource "aws_iam_policy" "dynamodb" {
  name        = "lambda_${local.name}_dynamodb"
  path        = "/"
  description = "IAM policy to allow lambda to access to Dynamo DB"

  policy = data.aws_iam_policy_document.dynamodb.json
}

resource "aws_dynamodb_table" "dynamodb" {
  name = local.name
  billing_mode = "PAY_PER_REQUEST"
  hash_key = "id"

  attribute {
    name = "id"
    type = "S"
  }

  ttl {
    attribute_name = "time_to_exist"
    enabled = true
  }
}  

resource "aws_iam_role_policy_attachment" "dynamodb" {
  role       = aws_iam_role.iam_role.name
  policy_arn = aws_iam_policy.dynamodb.arn
}
