resource "aws_ecs_cluster" "fargate_cluster" {
  name = "${var.cluster_prefix}-${var.aws_region}-fargate"

  setting {
    name  = "containerInsights"
    value = "enabled"
  }
}

resource "aws_ecs_cluster_capacity_providers" "capacity_provider" {
  cluster_name = aws_ecs_cluster.cluster.name

  capacity_providers = ["FARGATE"]

  default_capacity_provider_strategy {
    base              = 1
    weight            = 100
    capacity_provider = "FARGATE"
  }
}

resource "aws_cloudwatch_log_group" "main" {
  name              = "/aws/ecs/${var.cluster_name}/${var.service_name}/${var.environment}"
  retention_in_days = 3
}