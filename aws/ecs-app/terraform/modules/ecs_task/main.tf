data "aws_subnets" "public_subnets" {
  filter {
    name   = "tag:Name"
    values = ["*-pub-eu-west-1"]
  }
}

resource "aws_ecs_service" "service" {
  name             = var.service_name
  cluster          = data.aws_ecs_cluster.dzbt_bi_analytics.id
  task_definition  = aws_ecs_task_definition.service.arn
  platform_version = "1.4.0"
  desired_count    = var.desired_count

  network_configuration {
    subnets          = data.aws_subnets.public_subnets.ids
    security_groups  = [aws_security_group.ecs.id]
    assign_public_ip = false
  }

  capacity_provider_strategy {
    capacity_provider = "FARGATE"
    weight            = 100
  }
}

resource "aws_ecs_task_definition" "service" {
  family = var.service_name

  container_definitions = templatefile("${path.module}/templates/task-container-definitions-template.tftpl", {
    docker_image_url      = var.docker_image_url
    service_name          = var.service_name
    cloudwatch_log_region = var.aws_region
    cloudwatch_log_group  = aws_cloudwatch_log_group.main.name
    container_port        = var.container_port
    cpu                   = var.cpu
    memory_reservation    = var.memory
    environment = jsonencode([
      { name : "ENVIRONMENT", value : var.environment },
    ])
  })
  requires_compatibilities = ["FARGATE"]
  cpu                      = var.cpu
  memory                   = var.memory
  task_role_arn            = aws_iam_role.ecs_task_execution.arn
  execution_role_arn       = aws_iam_role.ecs_task_execution.arn
  network_mode             = "awsvpc"
  runtime_platform {
    cpu_architecture = "ARM64"
  }
}

