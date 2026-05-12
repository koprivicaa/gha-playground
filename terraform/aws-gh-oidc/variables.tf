variable "aws_region" {
  description = "AWS region for provider operations"
  type        = string
  default     = "eu-central-1"
}

variable "github_oidc_url" {
  description = "OIDC issuer URL used by GitHub Actions"
  type        = string
  default     = "https://token.actions.githubusercontent.com"
}

variable "github_oidc_client_ids" {
  description = "Allowed client IDs (audiences) for the OIDC provider"
  type        = list(string)
  default     = ["sts.amazonaws.com"]
}

variable "github_oidc_audience" {
  description = "Audience claim required in GitHub OIDC token when assuming the role"
  type        = string
  default     = "sts.amazonaws.com"
}

variable "github_repository" {
  description = "GitHub repository in owner/name format allowed to assume the role"
  type        = string
  default     = "koprivicaa/gha-playground"
}

variable "github_branch" {
  description = "Git branch allowed to assume the role"
  type        = string
  default     = "main"
}

variable "github_actions_role_name" {
  description = "Name of IAM role assumed by GitHub Actions"
  type        = string
  default     = "gha-playground-main-role"
}

variable "github_actions_policy_name" {
  description = "Name of IAM policy attached to the GitHub Actions role"
  type        = string
  default     = "gha-playground-github-actions-permissions"
}

variable "github_actions_allowed_actions" {
  description = "Allowed AWS API actions for the GitHub Actions role"
  type        = list(string)
  default = [
    "sts:GetCallerIdentity",
    "s3:ListAllMyBuckets",
  ]
}

variable "existing_github_actions_role_name" {
  description = "Pre-existing IAM role name used by GitHub Actions workflows"
  type        = string
  default     = "github-actions-role"
}

variable "deploy_s3_bucket_name" {
  description = "S3 bucket name used by deploy workflow upload step"
  type        = string
  default     = "gha-playground-demo"
}

variable "deploy_s3_object_key" {
  description = "S3 object key uploaded by deploy workflow"
  type        = string
  default     = "demo.txt"
}
