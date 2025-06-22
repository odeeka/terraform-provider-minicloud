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

# data "minicloud_vms" "demo" {}

# output "vms" {
#   value = data.minicloud_vms.demo
# }

# resource "minicloud_vms" "test" {
#   name   = "tf-vm"
#   image  = "nginx:latest"
#   cpu    = 0.5
#   memory = 80
#   # env = {
#   #   "ENV" = "Demo"
#   # }
#   ports = [6666]
# }

# resource "minicloud_vms" "test_2" {
#   name   = "tf-vm-2"
#   image  = "nginx:alpine"
#   cpu    = 1.5
#   memory = 88
#   ports  = [1000, 2000]
#   env = {
#     "ENV"       = "Development"
#     "Region"    = "WEU"
#     "Country"   = "CH"
#     "Continent" = "Europe"
#   }
# }

# resource "minicloud_vms" "test_3" {
#   name   = "tf-vm-3"
#   image  = "nginx:alpine"
#   cpu    = 1.0
#   memory = 999
#   ports  = [1100, 2200, 3300, 4400]
#   env = {
#     "ENV"     = "Staging"
#     "Region"  = "CHN"
#     "Country" = "CH"
#   }
# }

# resource "minicloud_vms" "test_2" {
#   name   = "tf-vm-2"
#   image  = "nginx:latest"
#   cpu    = 0.5
#   memory = 64
# }

locals {
  vms = {
    # tf-vm-a = {
    #   memory = 64
    # }
    test-vm = {
      memory = 64
    }
  }
}

# resource "minicloud_vms" "multi_vm" {
#   for_each = local.vms
#   name     = each.key
#   image    = "nginx:latest"
#   cpu      = 0.5
#   memory   = each.value.memory
# }

# resource "minicloud_vms" "multi" {
#   for_each = local.vms
#   name     = each.key
#   image    = "nginx:latest"
#   cpu      = 0.5
#   memory   = each.value.memory
# }