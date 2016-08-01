variable "do_token" {}
variable "pub_key" {}
variable "pvt_key" {}
variable "ssh_fingerprint" {}

provider "digitalocean" {
  token = "${var.do_token}"
}

resource "digitalocean_droplet" "www-1" {
    image = "ubuntu-16-04-x64"
    name = "www-1"
    region = "fra1"
    size = "512mb"
    private_networking = true
    ssh_keys = [
      "${var.ssh_fingerprint}"
    ]

    connection {
        user = "root"
        type = "ssh"
        key_file = "${var.pvt_key}"
        timeout = "2m"
    }

    provisioner "remote-exec" {
       inline = [
         "curl -sSL https://get.docker.com/ | sh"
       ]
     }
}
