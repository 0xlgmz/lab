Absolutely — here’s a **realistic Terraform coding task** that you might get during your 1-hour live DevOps interview. It reflects real-world infrastructure, is scoped for about 30–40 minutes of hands-on work, and includes a challenge around **ECS on Fargate**, **Terraform modules**, **IAM**, and **networking**.

---

## 🔧 **Mock Terraform Task: Deploy a Web App on ECS with Fargate**

### 🔹 **Scenario:**

Your team is launching a simple containerized web app (Nginx) on **AWS ECS with Fargate**. You need to write Terraform to deploy the infrastructure in a **reusable way**, using modules when appropriate.

---

### ✅ **Requirements:**

1. **ECS Fargate setup:**

   * Deploy an ECS Cluster.
   * Create a task definition running the `nginx:latest` container on port 80.
   * ECS Service should run 2 tasks behind a Load Balancer.

2. **Networking:**

   * Create or use an existing VPC with **public subnets** (assume 2 AZs).
   * Create a **security group** allowing inbound HTTP (port 80) and egress all.

3. **Load Balancing:**

   * Create an **Application Load Balancer (ALB)** to distribute traffic to the ECS service.
   * Add an ALB listener for HTTP traffic on port 80.

4. **IAM:**

   * ECS task must run with a task execution role allowing pulling images from ECR and writing logs to CloudWatch.

5. **Logging:**

   * Enable container logging to **CloudWatch Logs**.

---

### 💡 **Bonus (if time allows):**

* Use a **Terraform module** for the ECS service or VPC setup.
* Output the ALB DNS name.

---

### 📝 **Instructions:**

* You can assume Terraform 1.4+ and AWS provider `~> 5.0`.
* Start from scratch or use minimal stubs.
* Aim for modular, readable, and well-organized Terraform code.
* Don’t worry about exact VPC CIDRs or AZ names — use data sources where appropriate.

---

## 🧪 Sample Directory Structure (Recommended)

```bash
ecs-app/
├── main.tf
├── variables.tf
├── outputs.tf
├── modules/
│   └── ecs_service/
│       ├── main.tf
│       ├── variables.tf
│       ├── outputs.tf
```

---

## 📦 Suggested Components to Include

### In `main.tf`:

* Provider block
* VPC (use `aws_vpc` or `terraform-aws-modules/vpc/aws`)
* ECS Cluster
* ALB + listener
* ECS module block (calls your `modules/ecs_service`)

### In `modules/ecs_service/`:

* Task definition
* Service
* IAM role for execution
* Log group
* Security group (if module-specific)

