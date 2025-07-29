variable "service_name" {
  description = "Name of the ECS service"
  type        = string
  default     = "ecs-service"
}
variable "docker_image_tag" {
  description = "Docker image tag for the ECS service"
  type        = string
  default     = "latest"
}
variable "container_port" {
  description = "Port on which the container listens"
  type        = number
  default     = 80
}
variable "cpu" {
  description = "CPU units for the ECS task"
  type        = number
  default     = 256
}
variable "memory" {
  description = "Memory in MiB for the ECS task"
  type        = number
  default     = 512
}
variable "desired_count" {
  description = "Desired number of ECS tasks"
  type        = number
  default     = 1
}
variable "docker_image_url" {
  description = "Docker image URL for the ECS task"
  type        = string
  default     = "nginx/nginx:${var.docker_image_tag}"
}
variable "aws_region" {
  description = "AWS region for the ECS task"
  type        = string
  default     = "eu-west-1" 
}
variable "environment" {
  description = "Environment for the ECS task (e.g., dev, staging, prod)"
  type        = string
  default     = "dev"
}
variable "cluster_prefix" {
  description = "Prefix for ECS cluster resources"
  type        = string
  default     = "ecs-cluster"
}