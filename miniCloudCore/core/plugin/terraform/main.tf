
terraform {
  required_providers {
    virtualbox = {
      source  = "shekeriev/virtualbox"
      version = "0.0.4"
    }
  }
}

provider "virtualbox" {
  delay      = 60
  mintimeout = 5
}

# There are currently no configuration options for the provider itself.
resource "virtualbox_vm" "node" {
  count     = 1
  name      = format("node-%02d", count.index + 1)
  image     = "https://app.vagrantup.com/shekeriev/boxes/debian-11/versions/0.2/providers/virtualbox.box"
  cpus      = 1
  memory    = "512 mib"
  user_data = file("user_data")


  network_adapter {
    type           = "bridged"
    host_interface = "enp0s3"
  }


}

output "IPAddr" {
  value = element(virtualbox_vm.node.*.network_adapter.0.ipv4_address, 1)
}

output "IPAddr_2" {
  value = element(virtualbox_vm.node.*.network_adapter.0.ipv4_address, 2)
}



 