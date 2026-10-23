output "bucket_id" {
  description = "The ID of the S3 bucket for card images."
  value       = aws_s3_bucket.card_image.id
}

output "bucket_name" {
  description = "The name of the S3 bucket for card images."
  value       = aws_s3_bucket.card_image.bucket
}

output "bucket_arn" {
  description = "The ARN of the S3 bucket for card images."
  value       = aws_s3_bucket.card_image.arn
}
