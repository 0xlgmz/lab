variable "cluster_prefix" {
  description = "Prefix for the ECS cluster name"
  type        = string
  default     = "shiny-cluster"
}
variable "aws_region" {
  description = "AWS region for the ECS cluster"
  type        = string
  default     = "eu-west-1"
}
variable "environment" {
  description = "Environment for the ECS cluster (e.g., dev, staging, prod)"
  type        = string
  default     = "dev"
}  
variable "cluster_name" {
  description = "Name of the ECS cluster"
  type        = string
  default     = "fargate-cluster"
}
variable "service_name" {
  description = "Name of the ECS service"
  type        = string
  default     = "fargate-service"
}