module "vpc_base" {
  source     = "./terraform_modules/vpc_base"
  cidr_block = "10.0.0.0/16"
  azs        = ["eu-west-1a", "eu-west-1b"]
  name       = "ec2_vpc"
  service_name = "ec2_application"
  environment = var.environment
  cluster_prefix = var.cluster_prefix
}

module "fargate_cluster" {
  source          = "./terraform/modules/fargate_cluster"
  cluster_prefix  = var.cluster_prefix
  aws_region      = var.aws_region
  environment     = var.environment
  cluster_name    = var.cluster_name
  service_name    = var.service_name
}

module "fargate_service" {
  source          = "./terraform/modules/ecs_task"
  service_name    = local.service_name
  docker_image_url = "nginx/nginx:${var.docker_image_tag}"
  container_port  = 80
  cpu             = var.cpu
  memory          = var.memory
  desired_count   = 2
  aws_region      = var.aws_region
  environment     = var.environment
  cluster_prefix  = var.cluster_prefix
  cluster_name    = local.cluster_name
  task_execution_role = aws_iam_role.ecs_task_execution_role.arn
}

// Application Load Balancer for the ECS service listening on port 80
resource "aws_lb" "ecs_alb" {
  name               = "${var.cluster_prefix}-ecs-alb"
  internal           = false
  load_balancer_type = "application"
  security_groups    = [aws_security_group.ecs_alb.id]
  subnets            = module.vpc_base.public_subnet_ids

  enable_deletion_protection = false

  tags = {
    Name        = "${var.cluster_prefix}-ecs-alb"
    Environment = var.environment
  }
}