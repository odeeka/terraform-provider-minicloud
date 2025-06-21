terraform {
  required_providers {
    minicloud = {
      source = "hashicorp.com/edu/minicloud"
    }
  }
}

provider "minicloud" {
  host     = "http://localhost:8080"
  username = "testuser"
  password = "TestPass2025"
}

data "minicloud_vms" "demo" {}

output "vms" {
  value = data.minicloud_vms.demo
}
