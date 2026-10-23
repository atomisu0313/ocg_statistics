output "lambda_arn" {
  description = "The ARN of the Lambda function."
  value       = aws_lambda_function.lambda_syncstorageimage.arn
}

output "lambda_name" {
  description = "The name of the Lambda function."
  value       = aws_lambda_function.lambda_syncstorageimage.function_name
}
