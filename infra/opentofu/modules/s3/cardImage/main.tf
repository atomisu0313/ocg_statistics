resource "aws_s3_bucket" "card_image" {
  bucket = var.bucket_name

  tags = {
    project = var.tag_base
  }
}

resource "aws_s3_bucket_public_access_block" "card_image" {
  bucket                  = aws_s3_bucket.card_image.id
  block_public_acls       = true
  block_public_policy     = true
  ignore_public_acls      = true
  restrict_public_buckets = true
}

resource "aws_s3_bucket_server_side_encryption_configuration" "card_image" {
  bucket = aws_s3_bucket.card_image.id

  rule {
    apply_server_side_encryption_by_default {
      sse_algorithm = "AES256"
    }
  }
}

resource "aws_s3_bucket_lifecycle_configuration" "card_image" {
  bucket = aws_s3_bucket.card_image.id

  rule {
    id     = "abort-incomplete-multipart-upload"
    status = "Enabled"

    filter {}

    abort_incomplete_multipart_upload {
      days_after_initiation = 7
    }
  }
}
