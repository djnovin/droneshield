# Output for the load balancer DNS name
output "elb_dns_name" {
  value = aws_lb.example.dns_name
}
