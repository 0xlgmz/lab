# Merge locals from both fargate and service modules here

locals {
  cluster_prefixes = {
    dev   = "crimson-otter"
    stage = "silent-moon"
    prod  = "amber-forest"
  }
  cluster_prefix = local.cluster_prefixes[var.environment]
  cluster_name   = "${local.cluster_prefix}-${var.aws_region}_${var.environment}_${var.aws_account}"

  service_name             = "dzbt-bi-analytics-${var.environment}"
  vpc_id                   = var.vpc_id
  cpu                      = 512
  memory                   = 1024
  bi_analytics_bucket_name = "${var.aws_account}-bi-analytics-${var.aws_region}"
}
