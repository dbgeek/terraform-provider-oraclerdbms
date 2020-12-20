terraform {
  required_providers {
    oraclerdbms = {
      source  = "local/tf/oraclerdbms"
      version = "~> 0.1"
    }
  }
  required_version = ">= 0.13"
}

resource "oraclerdbms_profile" "app_profile" {
  profile = "APP_PROFILE"
}
