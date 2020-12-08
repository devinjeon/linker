terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 3.0"
    }
    cloudflare = {
      source  = "cloudflare/cloudflare"
      version = "~> 2.0"
    }
  }
  backend "s3" {
    bucket = "devin-terraform"
    key    = "states/linker.tfstate"
    region = "ap-northeast-1"
  }
}

locals {
  region = "ap-northeast-2"
}

provider "aws" {
  region = local.region
}

provider "cloudflare" {}
