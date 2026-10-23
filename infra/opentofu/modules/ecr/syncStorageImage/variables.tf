variable "repo_name" {
  description = "The name of the ECR repository for syncStorageImage."
  type        = string
}

variable "tag_base" {
  description = "The base tag for resources."
  type        = string
}
