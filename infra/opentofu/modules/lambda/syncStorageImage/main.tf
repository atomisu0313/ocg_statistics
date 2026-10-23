resource "aws_lambda_function" "lambda_syncstorageimage" {
  function_name = var.lambda_name
  package_type  = "Image"
  role          = var.role_arn
  image_uri     = var.image_uri

  memory_size   = 1024
  timeout       = 300
  architectures = ["arm64"]

  ephemeral_storage {
    size = 512
  }

  environment {
    variables = var.environment
  }

  tags = {
    project = var.tag_base
  }
}
