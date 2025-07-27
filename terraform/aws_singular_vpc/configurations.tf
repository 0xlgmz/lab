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
variable "aws_region" {}
variable "managed" {}
variable "application_name" {}
variable "owner" {}
variable "repo" {}
variable "environment" {}
