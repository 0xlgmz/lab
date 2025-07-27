module "vpc_base" {
  source     = "./terraform_modules/vpc_base"
  cidr_block = "10.0.0.0/16"
  azs        = ["eu-west-1a", "eu-west-1b"]
  name       = "lab"
  service_name = "my_service"
  environment = var.environment
}