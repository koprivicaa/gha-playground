provider "aws" {
  region                      = var.aws_region
  skip_credentials_validation = true
  skip_metadata_api_check     = true
  skip_requesting_account_id  = true
}

data "tls_certificate" "github_actions" {
  url = var.github_oidc_url
}

resource "aws_iam_openid_connect_provider" "github" {
  url = var.github_oidc_url

  client_id_list = [var.github_oidc_audience]

  # Use the current SHA-1 fingerprint from GitHub's OIDC TLS certificate chain.
  thumbprint_list = [data.tls_certificate.github_actions.certificates[0].sha1_fingerprint]
}

data "aws_iam_policy_document" "github_actions_assume_role" {
  statement {
    effect  = "Allow"
    actions = ["sts:AssumeRoleWithWebIdentity"]

    principals {
      type        = "Federated"
      identifiers = [aws_iam_openid_connect_provider.github.arn]
    }

    condition {
      test     = "StringEquals"
      variable = "token.actions.githubusercontent.com:aud"
      values   = [var.github_oidc_audience]
    }

    condition {
      test     = "StringEquals"
      variable = "token.actions.githubusercontent.com:sub"
      values   = ["repo:${var.github_repository}:ref:refs/heads/${var.github_branch}"]
    }
  }
}

resource "aws_iam_role" "github_actions" {
  name               = var.github_actions_role_name
  assume_role_policy = data.aws_iam_policy_document.github_actions_assume_role.json
}

data "aws_iam_policy_document" "github_actions_permissions" {
  statement {
    effect    = "Allow"
    actions   = var.github_actions_allowed_actions
    resources = ["*"]
  }
}

resource "aws_iam_policy" "github_actions_permissions" {
  name   = var.github_actions_policy_name
  policy = data.aws_iam_policy_document.github_actions_permissions.json
}

resource "aws_iam_role_policy_attachment" "github_actions_permissions" {
  role       = aws_iam_role.github_actions.name
  policy_arn = aws_iam_policy.github_actions_permissions.arn
}

data "aws_iam_policy_document" "deploy_upload_permissions" {
  statement {
    effect    = "Allow"
    actions   = ["s3:PutObject"]
    resources = ["arn:aws:s3:::${var.deploy_s3_bucket_name}/${var.deploy_s3_object_key}"]
  }
}

resource "aws_iam_role_policy" "github_actions_deploy_upload" {
  name   = "github-actions-deploy-upload"
  role   = aws_iam_role.github_actions.name
  policy = data.aws_iam_policy_document.deploy_upload_permissions.json
}
