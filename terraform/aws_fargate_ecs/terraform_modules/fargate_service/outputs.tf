# Merge outputs from both fargate and service modules here

output "fargate_cluster_name" {
  value = aws_ecs_cluster.cluster.name
}

output "private_subnets" {
  value       = data.aws_subnets.private_subnets.ids
  description = "List of private subnets in the VPC"
}
