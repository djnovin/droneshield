terraform {
  required_version = ">= 1.0"

  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 5.0"
    }
  }
}

# Data source for available AZs in the region
data "aws_availability_zones" "available" {}

# Security group for both instances and ELB
resource "aws_security_group" "default" {
  name = "${var.PROJECT_NAME}-${var.ENVIRONMENT}-sg"

  ingress {
    from_port   = 80
    to_port     = 80
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }

  tags = {
    Name        = "${var.PROJECT_NAME}-${var.ENVIRONMENT}-sg"
    Environment = var.ENVIRONMENT
  }
}

# Launch configuration for the autoscaling group
resource "aws_launch_configuration" "example" {
  image_id        = "ami-58d7e821"
  instance_type   = "t2.micro"
  security_groups = [aws_security_group.default.id]

  user_data = <<-EOF
              #!/bin/bash
              echo "Hello, World" > index.html
              nohup busybox httpd -f -p 80 &
              EOF

  lifecycle {
    create_before_destroy = true
  }
}

# Autoscaling group for high availability
resource "aws_autoscaling_group" "example" {
  launch_configuration = aws_launch_configuration.example.id
  availability_zones   = data.aws_availability_zones.available.names
  min_size             = 2
  max_size             = 5
  desired_capacity     = 3

  target_group_arns = [aws_lb_target_group.example.arn]

  tag {
    key                 = "Name"
    value               = "${var.PROJECT_NAME}-${var.ENVIRONMENT}-asg"
    propagate_at_launch = true
  }
}

# Load balancer for distributing traffic
resource "aws_lb" "example" {
  name                       = "${var.PROJECT_NAME}-${var.ENVIRONMENT}-lb"
  internal                   = false
  load_balancer_type         = "application"
  security_groups            = [aws_security_group.default.id]
  subnets                    = data.aws_availability_zones.available.names
  enable_deletion_protection = false

  tags = {
    Name        = "${var.PROJECT_NAME}-${var.ENVIRONMENT}-lb"
    Environment = var.ENVIRONMENT
  }
}

# Target group for the load balancer
resource "aws_lb_target_group" "example" {
  name        = "${var.PROJECT_NAME}-${var.ENVIRONMENT}-tg"
  target_type = "instance"
  port        = 80
  protocol    = "HTTP"
  vpc_id      = var.VPC_ID

  health_check {
    interval            = 30
    timeout             = 5
    healthy_threshold   = 2
    unhealthy_threshold = 2
    path                = "/"
    protocol            = "HTTP"
  }

  tags = {
    Name        = "${var.PROJECT_NAME}-${var.ENVIRONMENT}-tg"
    Environment = var.ENVIRONMENT
  }
}
