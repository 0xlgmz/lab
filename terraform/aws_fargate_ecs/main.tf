module "vpc_base" {
  source     = "./terraform_modules/vpc_base"
  cidr_block = "10.0.0.0/16"
  azs        = ["eu-west-1a", "eu-west-1b"]
  name       = "lab"
  service_name = "my_service"
  environment = var.environment
}

module "merged_fargate_service" {
  source         = "/terraform_modules/fargate-service"
  aws_region     = var.aws_region
  aws_account    = var.aws_account
  environment    = var.environment
  desired_count  = 1
  docker_image_tag = "latest"
  cpu            = 512
  memory         = 1024
  container_port = 80
  aws_org_id     = "o-52d4n4n0uq"
  vpc_id         = module.vpc_base.vpc_id
  # Add other variables as needed
}