# Copy resources from fargate/cluster.tf, service/service.tf, service/iam.tf, service/s3.tf, service/eventbridge.tf, service/data.tf
# Example:
# resource "aws_ecs_cluster" "main" { ... }
# resource "aws_iam_role" "ecs_task_execution" { ... }
# resource "aws_s3_bucket" "service_data" { ... }
# resource "aws_eventbridge_rule" "service_event" { ... }

resource "aws_ecs_cluster" "cluster" {
  name = "${local.cluster_prefix}-${var.aws_region}-fargate"

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
  name              = "/aws/ecs/${local.cluster_name}/${local.service_name}/${var.environment}"
  retention_in_days = 3
}

data "aws_subnets" "private_subnets" {
  filter {
    name   = "tag:Name"
    values = ["*-pvt-eu-west-1"]
  }
}

data "aws_ecs_cluster" "dzbt_bi_analytics" {
  cluster_name = "${local.cluster_prefix}-${var.aws_region}-fargate"
}

resource "aws_ecs_service" "service" {
  name             = local.service_name
  cluster          = data.aws_ecs_cluster.dzbt_bi_analytics.id
  task_definition  = aws_ecs_task_definition.service.arn
  platform_version = "1.4.0"
  desired_count    = var.desired_count

  network_configuration {
    subnets          = data.aws_subnets.private_subnets.ids
    security_groups  = [aws_security_group.ecs.id]
    assign_public_ip = false
  }

  capacity_provider_strategy {
    capacity_provider = "FARGATE"
    weight            = 100
  }
}

resource "aws_ecs_task_definition" "service" {
  family = local.service_name

  container_definitions = templatefile("${path.module}/templates/task-container-definitions-template.tftpl", {
    docker_image_url      = "070643671424.dkr.ecr.eu-west-1.amazonaws.com/dzbt-bi-analytics:${var.docker_image_tag}-emulated"
    service_name          = local.service_name
    cloudwatch_log_region = var.aws_region
    cloudwatch_log_group  = aws_cloudwatch_log_group.main.name
    container_port        = var.container_port
    cpu                   = local.cpu
    memory_reservation    = local.memory
    environment = jsonencode([
      { name : "ENVIRONMENT", value : var.environment },
    ])
  })
  requires_compatibilities = ["FARGATE"]
  cpu                      = local.cpu
  memory                   = local.memory
  task_role_arn            = aws_iam_role.ecs_task_execution.arn
  execution_role_arn       = aws_iam_role.ecs_task_execution.arn
  network_mode             = "awsvpc"
  runtime_platform {
    cpu_architecture = "ARM64"
  }
}

# Add resources from iam.tf, s3.tf, eventbridge.tf, data.tf here as needed
