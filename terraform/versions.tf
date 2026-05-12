terraform {
  required_version = ">= 1.6.0"
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 5.0"
    }
  }
}

#Napiši Terraform za:
#1. `aws_iam_openid_connect_provider` za GitHub
#2. `aws_iam_role` sa trust policy koji ograničava na `repo:koprivicaa/gha-playground:ref:refs/heads/main`
#3. Policy koja dozvoljava samo: `sts:GetCallerIdentity`, `s3:ListAllMyBuckets`, 's3:GetObject' i 's3:PutObject' na bucketu `gha-playground-deployments`
#4. Output role ARN

#Apply, zabeleži ARN.