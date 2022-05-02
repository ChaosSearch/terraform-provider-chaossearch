/*  
  All CS prefix vars have a matching EnvDefaultFunc
  Ex, 'CS_URL' == export CS_URL='example.chaossearch.io'
  
  Or if you'd like to use vars.tf you can export as TF_VAR_CS_URL
  Or add them to your .tfvars
*/

variable "CS_USERNAME" {
  type = string
  description = "Username"
}

variable "CS_PASSWORD" {
  type = string
  description = "Password"
}

variable "CS_PARENT_USER_ID" {
  type = string
  description = "Parent User Id"
}

variable "CS_URL" {
  type = string
  description = "Cluster Url"
}

variable "CS_ACCESS_KEY" {
  type = string
  description = "Cluster Access Key Id"
}

variable "CS_SECRET_KEY" {
  type = string
  description = "Cluster Secret Access Key"
}

variable "CS_REGION" {
  type = string
  description = "Cluster Region"
}

variable "AWS_PROFILE" {
  type = string
  description = "AWS profile"
}