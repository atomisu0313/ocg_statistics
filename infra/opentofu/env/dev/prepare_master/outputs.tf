output "s3_card_image_bucket_name" {
  description = "The name of the S3 bucket for card images."
  value       = module.s3_card_image.bucket_name
}

output "s3_card_image_bucket_arn" {
  description = "The ARN of the S3 bucket for card images."
  value       = module.s3_card_image.bucket_arn
}

output "lambda_syncstorageimage_arn" {
  description = "The ARN of the syncStorageImage Lambda function."
  value       = module.lambda_syncstorageimage.lambda_arn
}

output "lambda_syncstorageimage_name" {
  description = "The name of the syncStorageImage Lambda function."
  value       = module.lambda_syncstorageimage.lambda_name
}

