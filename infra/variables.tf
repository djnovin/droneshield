variable "REGION" {
  description = "AWS region"
  type        = string
  default     = "ap-southeast-2"
}

# Variables for dynamic naming
variable "PROJECT_NAME" {
  description = "Name of the project"
  default     = "droneshield"
}

variable "ENVIRONMENT" {
  description = "Deployment environment (e.g., dev, staging, prod)"
  default     = "dev"
}

variable "VPC_ID" {
  description = "ID of the VPC"
  type        = string
}
