# terraform-provider-edgeos

Terraform provider for managing EdgeOS configs. Only works with firewall rules on an existing firewall.

Example

```terraform
terraform {
  required_providers {
    edgeos = {
      version = "~> 0.0.1"
      source  = "montanadev.com/pkg/edgeos"
    }
  }
}

provider "edgeos" {
  username = var.username
  password = var.password
  host     = var.host
}

resource "edgeos_firewall_rule" "test_vlan" {
  firewall_name = "FIREWALL"
  rule_id       = "250"
  name          = "vlan access across 5051"
  action        = "accept"
  protocol      = "tcp_udp"

  source_address      = "10.0.23.0/24"
  destination_address = "10.0.24.28"
  destination_port    = 5051

  states = ["established", "related"]
}
```