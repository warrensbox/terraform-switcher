terraform {
  required_version = "~> 1.0.0"

  required_providers {
    aws        = ">= 2.52.0"
    kubernetes = ">= 1.11.1"
  }
}

terraform {
  required_version = "<= 1.0.5"
}
