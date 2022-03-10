variable "user_name" {
  type = string
  description = "Username"
}

variable "password" {
  type = string
  description = "Password"
}

variable "parent_user_id" {
  type = string
  description = "Parent User Id"
}

variable "admin_user_name" {
  type = string
  description = "Admin Username"
}

variable "admin_password" {
  type = string
  description = "Admin Password"
}

variable "url" {
  type = string
  description = "Cluster Url"
}

variable "access_key_id" {
  type = string
  description = "Cluster Access Key Id"
}

variable "secret_access_key" {
  type = string
  description = "Cluster Secret Access Key"
}

variable "region" {
  type = string
  description = "Cluster Region"
}

variable "profile" {
  type = string
  description = "aws profile"
}