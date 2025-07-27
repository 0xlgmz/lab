variable "environment" {
  description = "The environment for the ECS service (e.g., dev, staging, prod)"
  type        = string
}

variable "aws_region" {
  description = "The AWS region for the ECS cluster"
  type        = string
}

variable "aws_account" {
  description = "The AWS account ID"
  type        = string
}

variable "desired_count" {
  description = "The desired number of tasks for the ECS service"
  type        = number
  default     = 0
}

variable "docker_image_tag" {
  description = "The Docker image tag for the ECS service"
  type        = string
}

variable "cpu" {
  description = "The amount of CPU units to allocate for the ECS task"
  type        = number
  default     = 256
}

variable "memory" {
  description = "The amount of memory (in MiB) to allocate for the ECS task"
  type        = number
  default     = 512
}

variable "container_port" {
  description = "The port on which the container listens"
  type        = number
  default     = 80
}

variable "aws_org_id" {
  description = "organisation id"
  type        = string
  default     = "o-52d4n4n0uq"
}

variable "vpc_id" {
  description = "The ID of the VPC where the ECS service will be deployed"
  type        = string
}