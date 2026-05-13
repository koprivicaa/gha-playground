data "aws_caller_identity" "current" {}

variable "org_prefix" {
  type    = string
  default = "acme"
}

variable "environment" {
  type    = string
  default = "prod"
}

variable "cloudtrail_name" {
  type    = string
  default = "my-cloudtrail"
}

variable "cloudtrail_s3_key_prefix" {
  type    = string
  default = "cloudtrail"
}

variable "cloudtrail_sensitive_data_bucket_suffix" {
  type    = string
  default = "sensitive-data"
}

variable "cloudtrail_lambda_data_resource_arn" {
  type    = string
  default = "arn:aws:lambda:*"
}

variable "cloudtrail_tag_name" {
  type    = string
  default = "org-trail"
}

variable "cloudtrail_compliance_tag" {
  type    = string
  default = "SOC2,PCI-DSS"
}

variable "cloudwatch_log_group_name" {
  type    = string
  default = "/aws/cloudtrail/org-trail"
}

variable "cloudwatch_retention_days" {
  type    = number
  default = 365
}

variable "cloudtrail_to_cwlogs_role_name" {
  type    = string
  default = "cloudtrail-to-cwlogs"
}

variable "kms_key_description" {
  type    = string
  default = "CMK for CloudTrail log encryption"
}

variable "kms_deletion_window_in_days" {
  type    = number
  default = 30
}

variable "kms_key_tag_name" {
  type    = string
  default = "cloudtrail-cmk"
}

variable "kms_alias_name" {
  type    = string
  default = "alias/cloudtrail"
}

variable "cloudtrail_bucket_force_destroy" {
  type    = bool
  default = false
}

variable "cloudtrail_bucket_tag_name" {
  type    = string
  default = "cloudtrail-logs"
}

variable "cloudtrail_lifecycle_transition_days" {
  type    = number
  default = 90
}

variable "cloudtrail_lifecycle_storage_class" {
  type    = string
  default = "GLACIER"
}

variable "cloudtrail_lifecycle_expiration_days" {
  type    = number
  default = 2555
}

variable "cloudtrail_lifecycle_noncurrent_expiration_days" {
  type    = number
  default = 90
}