// Backend Configuration
terraform {
  required_providers {
    aws = {
        source = "hashicorp/aws"
        version = "= 1.12.2"
    }
  }
  backend "s3" {}
}

// Provider Configuration
provider "aws" {
  region = var.aws_region
  default_tags {
    tags = {
      ManagedBy = "This is managed by ${var.managed} and deployed via Terraform"
      Application = var.application_name
      Owner = var.owner
      GitRepo = var.repo
      Environment = var.environment
    }
  }
}

// Variables
variable "aws_region" {
  description = "AWS region where resources will be deployed"
  type        = string
  default = "eu-west-1"
}

variable "cluster_prefix" {
  description = "Prefix for ECS cluster names"
  type        = string
  default = "shiny-flash"
}

variable "environment" {
  description = "Deployment environment (e.g., dev, prod)"
  type        = string
  default     = "dev"
}
variable "managed" {
  description = "Indicates if the infrastructure is managed"
  type        = string
  default     = "Terraform"
}
variable "application_name" {
  description = "Name of the application"
  type        = string
  default     = "ecs-application"
}
variable "owner" {
  description = "Owner of the resources"
  type        = string
  default     = "0xlgmz"
}
variable "repo" {
  description = "Git repository URL"
  type        = string
  default     = "github.com/0xlgmz/ecs-app"
}

// Locals Variables
locals {
  cluster_name = "ecs-cluster-${var.environment}"
  service_name = "ecs-service-${var.environment}"
}


// Output Variables
output "fargate_cluster_name" {
  value = aws_ecs_cluster.cluster.name
}