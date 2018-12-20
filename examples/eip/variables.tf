variable "region" {
  description = "The region that will create resources in"
  default     = "cn-bj2"
}

variable "instance_password" {
  default = "wA123456"
}

variable "count" {
  default = 3
}

variable "count_format" {
  default = "%02d"
}
