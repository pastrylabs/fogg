# Auto-generated by fogg. Do not edit
# Make improvements in fogg, so that everyone can benefit.


# Default Provider
provider "aws" {
  version             = "~> 0.13.0"
  region              = "us-west-bar1"
  profile             = "czi-bar"
  allowed_account_ids = [3]
}

# Aliased Providers (for doing things in every region).

provider "aws" {
  alias               = "us-east-bar2"
  version             = "~> 0.13.0"
  region              = "us-east-bar2"
  profile             = "czi-bar"
  allowed_account_ids = [3]
}











terraform {
  required_version = "=0.13.0"

  backend "s3" {
    bucket         = "bar-bucket"
    dynamodb_table = "bar-table"
    key            = "terraform/bar-project/accounts/bar.tfstate"
    encrypt        = true
    region         = "us-west-bar1"
    profile        = "czi-bar"
  }
}

variable "project" {
  type    = string
  default = "bar-project"
}


variable "region" {
  type    = string
  default = "us-west-bar1"
}



variable "aws_profile" {
  type    = string
  default = "czi-bar"
}


variable "owner" {
  type    = string
  default = "bar@example.com"
}

variable "aws_accounts" {
  type = map
  default = {


    bar = 3



    foo = 2


  }
}


variable "foo" {
  type    = string
  default = "bar"
}


data "terraform_remote_state" "global" {
  backend = "s3"

  config = {
    bucket         = "the-bucket"
    dynamodb_table = "the-table"
    key            = "terraform/test-project/global.tfstate"
    region         = "us-west-2"
    profile        = "czi"
  }
}






data "terraform_remote_state" "foo" {
  backend = "s3"

  config = {
    bucket         = "bar-bucket"
    dynamodb_table = "bar-table"
    key            = "terraform/bar-project/accounts/foo.tfstate"
    region         = "us-west-bar1"
    profile        = "czi-bar"
  }
}


