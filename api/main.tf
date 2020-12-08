locals {
  name        = "linker"
  target_file = "deploy.zip"
}

resource "aws_lambda_function" "lambda" {
  filename      = data.archive_file.lambda_deployment.output_path
  function_name = local.name
  role          = aws_iam_role.iam_role.arn
  handler       = "linker"
  runtime       = "go1.x"

  source_code_hash = data.archive_file.lambda_deployment.output_base64sha256

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

data "archive_file" "lambda_deployment" {
  type        = "zip"
  source_file = "src/linker"
  output_path = "deploy.zip"
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
    actions = [
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
  name         = local.name
  billing_mode = "PAY_PER_REQUEST"
  hash_key     = "ID"

  attribute {
    name = "ID"
    type = "S"
  }

  ttl {
    attribute_name = "time_to_exist"
    enabled        = true
  }
}

resource "aws_iam_role_policy_attachment" "dynamodb" {
  role       = aws_iam_role.iam_role.name
  policy_arn = aws_iam_policy.dynamodb.arn
}

# === API GATEWAY ===
resource "aws_api_gateway_rest_api" "api" {
  name        = "linker"
  description = "linker"

  endpoint_configuration {
    types = ["REGIONAL"]
  }
}

resource "aws_api_gateway_resource" "api" {
  rest_api_id = aws_api_gateway_rest_api.api.id
  parent_id   = aws_api_gateway_rest_api.api.root_resource_id
  path_part   = "{id}"
}

resource "aws_api_gateway_method" "api" {
  rest_api_id   = aws_api_gateway_rest_api.api.id
  resource_id   = aws_api_gateway_resource.api.id
  http_method   = "ANY"
  authorization = "NONE"
}

resource "aws_api_gateway_integration" "api" {
  rest_api_id             = aws_api_gateway_rest_api.api.id
  resource_id             = aws_api_gateway_resource.api.id
  http_method             = aws_api_gateway_method.api.http_method
  integration_http_method = "POST"
  type                    = "AWS_PROXY"
  uri                     = aws_lambda_function.lambda.invoke_arn
}

resource "aws_lambda_permission" "api" {
  statement_id  = "AllowExecutionLinkerFromAPIGateway"
  action        = "lambda:InvokeFunction"
  function_name = aws_lambda_function.lambda.function_name
  principal     = "apigateway.amazonaws.com"

  source_arn = "${aws_api_gateway_rest_api.api.execution_arn}/*/*"
}

# === API GATEWAY - Deployment ===
resource "aws_api_gateway_deployment" "api" {
  depends_on = [
    aws_api_gateway_integration.api,
    aws_lambda_permission.api,
  ]

  stage_description = data.archive_file.lambda_deployment.output_base64sha256

  rest_api_id = aws_api_gateway_rest_api.api.id
  stage_name  = "main"
}

resource "aws_api_gateway_method_settings" "api" {
  depends_on = [
    aws_api_gateway_method.api,
    aws_api_gateway_deployment.api,
  ]

  rest_api_id = aws_api_gateway_rest_api.api.id
  stage_name  = aws_api_gateway_deployment.api.stage_name
  method_path = "*/*"

  settings {
    throttling_burst_limit = 200
    throttling_rate_limit  = 100
  }
}

variable "root_domain" {
  type = string
}

locals {
  custom_domain_name = "linker"
  custom_domain      = "${local.custom_domain_name}.${var.root_domain}"
}

resource "aws_acm_certificate" "domain" {
  domain_name       = local.custom_domain
  validation_method = "DNS"

  lifecycle {
    create_before_destroy = true
  }
}

resource "aws_acm_certificate_validation" "domain" {
  certificate_arn         = aws_acm_certificate.domain.arn
  validation_record_fqdns = [cloudflare_record.domain_validation.name]
}

resource "aws_api_gateway_domain_name" "domain" {
  domain_name              = local.custom_domain
  regional_certificate_arn = aws_acm_certificate_validation.domain.certificate_arn

  endpoint_configuration {
    types = ["REGIONAL"]
  }
}

resource "aws_api_gateway_base_path_mapping" "domain" {
  api_id      = aws_api_gateway_rest_api.api.id
  stage_name  = aws_api_gateway_deployment.api.stage_name
  domain_name = aws_api_gateway_domain_name.domain.domain_name
}

# === Cloudflare ===
locals {
  # NOTE(devin): `domain_validation_options` is 'set' type
  # (ref: https://github.com/hashicorp/terraform-provider-aws/issues/10098#issuecomment-663562342)
  validation         = tolist(aws_acm_certificate.domain.domain_validation_options)[0]
  cloudflare_zone_id = data.cloudflare_zones.domain.zones[0].id
}

data "cloudflare_zones" "domain" {
  filter {
    name = var.root_domain
  }
}

# 1. Create domain record to validate ACM certificate
resource "cloudflare_record" "domain_validation" {
  name    = local.validation.resource_record_name
  zone_id = local.cloudflare_zone_id

  # NOTE(devin): `resource_record_value` has suffix '.' but Cloudflare removes it automatically.
  value = trimsuffix(local.validation.resource_record_value, ".")
  type  = local.validation.resource_record_type
}

# 2. Create record -> links.{your_domain}
resource "cloudflare_record" "domain" {
  name    = local.custom_domain_name
  zone_id = local.cloudflare_zone_id

  value   = aws_api_gateway_domain_name.domain.regional_domain_name
  type    = "CNAME"
  proxied = true
}

# 3. Make page rules for forwarding url
resource "cloudflare_page_rule" "domain" {
  zone_id = local.cloudflare_zone_id
  target  = "${var.root_domain}/~*"

  actions {
    forwarding_url {
      url         = "https://${local.custom_domain}/$1"
      status_code = "301"
    }
  }
}
