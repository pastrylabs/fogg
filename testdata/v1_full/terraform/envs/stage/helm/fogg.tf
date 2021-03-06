# Auto-generated by fogg. Do not edit
# Make improvements in fogg, so that everyone can benefit.


provider "aws" {
  version             = "~> 0.14.0"
  region              = "us-west-stage1"
  profile             = "czi-stage"
  allowed_account_ids = [4]
}

# Aliased Providers (for doing things in every region).


provider "aws" {
  alias               = "us-east-stage2"
  version             = "~> 0.14.0"
  region              = "us-east-stage2"
  profile             = "czi-stage"
  allowed_account_ids = [4]
}












terraform {
  required_version = "~>0.14.0"

  backend "s3" {
    bucket         = "stage-bucket"
    dynamodb_table = "stage-table"

    key = "terraform/stage-project/envs/stage/components/helm.tfstate"


    encrypt = true
    region  = "us-west-stage1"
    profile = "czi-stage"
  }
}

variable "env" {
  type    = string
  default = "stage"
}

variable "project" {
  type    = string
  default = "stage-project"
}


variable "region" {
  type    = string
  default = "us-west-stage1"
}


variable "component" {
  type    = string
  default = "helm"
}


variable "aws_profile" {
  type    = string
  default = "czi-stage"
}



variable "owner" {
  type    = string
  default = "stage@example.com"
}

variable "tags" {
  type = map(string)
  default = {
    project   = "stage-project"
    env       = "stage"
    service   = "helm"
    owner     = "stage@example.com"
    managedBy = "terraform"
  }
}


variable "foo" {
  type    = string
  default = "stage"
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



data "terraform_remote_state" "cloud-env" {
  backend = "s3"

  config = {
    bucket         = "stage-bucket"
    dynamodb_table = "stage-table"
    key            = "terraform/stage-project/envs/stage/components/cloud-env.tfstate"
    region         = "us-west-stage1"
    profile        = "czi-stage"
  }
}


# remote state for accounts

data "terraform_remote_state" "bar" {
  backend = "s3"

  config = {
    bucket         = "bar-bucket"
    dynamodb_table = "bar-table"
    key            = "terraform/bar-project/accounts/bar.tfstate"
    region         = "us-west-bar1"
    profile        = "czi-bar"
  }
}

data "terraform_remote_state" "foo" {
  backend = "s3"

  config = {
    bucket         = "foo-bucket"
    dynamodb_table = "foo-table"
    key            = "terraform/foo-project/accounts/foo.tfstate"
    region         = "us-west-foo1"
    profile        = "czi-foo"
  }
}


# map of aws_accounts
variable "aws_accounts" {
  type = map
  default = {


    bar = 3



    foo = 2


  }
}
