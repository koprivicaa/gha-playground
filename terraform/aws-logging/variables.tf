variable "deploy_s3_bucket_name" {
  description = "S3 bucket name used by deploy workflow upload step"
  type        = string
  default     = "s3-cloudtrail"
}

data "aws_caller_identity" "current" {}

variable "org_prefix" {
  type    = string
  default = "acme"
}

variable "environment" {
  type    = string
  default = "prod"
}