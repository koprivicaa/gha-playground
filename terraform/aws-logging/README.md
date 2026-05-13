# Practical Task (60–90 min, hands-on)

Open a new Terraform layer in a test repo (it can also be local without AWS access, using only `terraform validate`).

---

## Goal

Create a minimal but **production-quality** logging setup for a CI/CD account.

---

## Resources to Implement

### 1. CloudTrail Organization Trail

```hcl
resource "aws_cloudtrail" "main" {
  ...
}
```

#### Requirements

- `is_organization_trail = true`
- `enable_log_file_validation = true`
- `kms_key_id = aws_kms_key.cloudtrail.arn`
- `include_global_service_events = true`
- S3 bucket with a policy allowing `cloudtrail.amazonaws.com` to write

---

### 2. CloudWatch Log Group for GHA Runners

```hcl
resource "aws_cloudwatch_log_group" "gha_runners" {
  ...
}
```

#### Requirements

- `retention_in_days = 30`
- `kms_key_id = aws_kms_key.logs.arn`

---

### 3. Kinesis Firehose → Splunk HEC

```hcl
resource "aws_kinesis_firehose_delivery_stream" "splunk" {
  ...
}
```

#### Requirements

- Destination: `splunk`
- Backup S3 bucket with **minimum 7-day retention**
- Lambda processor (placeholder ARN only)

---

### 4. Subscription Filter

Connect the CloudWatch log group to Firehose.

```hcl
resource "aws_cloudwatch_log_subscription_filter" "to_firehose" {
  ...
}
```

---

### 5. VPC Flow Logs

Enable Flow Logs for the VPC where the runners are hosted.

```hcl
resource "aws_flow_log" "vpc" {
  ...
}
```

#### Requirements

- `log_destination_type = "s3"`
- `traffic_type = "ALL"`
- Custom format string including all available fields

---

### 6. KMS Key

Create a KMS key with a policy that:

- Allows `logs.${region}.amazonaws.com` to use the key for encrypt/decrypt
- Allows `cloudtrail.amazonaws.com` to use the key
- Forbids direct root account access  
  (best practice — all access should go through IAM roles)

---

## Acceptance Criteria

- `terraform validate` passes
- `terraform plan` shows:
  - encryption enabled
  - retention configured
  - policies attached

---

## Required Outputs (`outputs.tf`)

Export the following values:

```hcl
output "log_group_arn" {}
output "firehose_arn" {}
output "cloudtrail_bucket_arn" {}
output "kms_key_arn" {}
```