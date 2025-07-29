variable "cidr_block" {
  description = "VPC CIDR block"
  type        = string
}

variable "azs" {
  description = "Availability zones"
  type        = list(string)
}

variable "name" {
  description = "Base name for resources"
  type        = string
}

variable "environment" {
  description = "Deployment environment (e.g., dev, prod)"
  type        = string
  default = "dev"
}

variable "service_name" {
  description = "Name of the service for labeling resources"
  type        = string
}

variable "cluster_prefix" {
  description = "Prefix for ECS cluster resources"
  type        = string
  default     = "ecs-cluster" 
}