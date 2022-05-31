# OKTA_ORG_NAME, OKTA_BASE_URL, and OKTA_API_TOKEN provided by env variables
# https://registry.terraform.io/providers/okta/okta/latest/docs
provider "okta" {}

terraform {
  required_providers {
    okta = {
      source = "okta/okta"
      version = "~> 3.20"
    }
  }
}