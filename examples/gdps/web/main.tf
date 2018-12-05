# Specify the provider and access details
provider "ucloud" {
  region = "${var.region}"
}

# Query availability zone
data "ucloud_zones" "default" {}

# Bulid instance type
data "ucloud_instance_types" "default" {
  cpu    = 1
  memory = 1
}

# Query image
data "ucloud_images" "default" {
  availability_zone = "${data.ucloud_zones.default.zones.0.id}"
  name_regex        = "^CentOS 7.[1-2] 64"
  image_type        = "Base"
}

# Create security group
resource "ucloud_security_group" "default" {
  name = "tf-example-gdps-web"
  tag  = "tf-example"

  # HTTP access from LAN
  rules {
    port_range = "80"
    protocol   = "TCP"
    cidr_block = "${var.subnet_network}"
    policy     = "ACCEPT"
  }

  # HTTPS access from LAN
  rules {
    port_range = "443"
    protocol   = "TCP"
    cidr_block = "${var.subnet_network}"
    policy     = "ACCEPT"
  }
}

# Create vpc
resource "ucloud_vpc" "default" {
  name = "tf-example-gdps-web"
  tag  = "tf-example"

  # vpc network
  cidr_blocks = ["${var.vpc_network}"]
}

# Create subnet
resource "ucloud_subnet" "default" {
  name = "tf-example-gdps-web"
  tag  = "tf-example"

  # subnet's network must be contained by vpc network
  # and a subnet must have least 8 ip addresses in it (netmask < 30).
  cidr_block = "${var.subnet_network}"

  vpc_id = "${ucloud_vpc.default.id}"
}

# Random shuffle will random select a zone
resource "random_shuffle" "az" {
  input        = ["${data.ucloud_zones.default.ids}"]
  result_count = "${var.instance_count}"
}

# Create a web server
resource "ucloud_instance" "web" {
  name              = "tf-example-gdps-web-${format(var.count_format, count.index+1)}"
  tag               = "tf-example"
  availability_zone = "${element(random_shuffle.az.result, count.index)}"
  image_id          = "${data.ucloud_images.default.images.0.id}"
  instance_type     = "${data.ucloud_instance_types.default.instance_types.0.id}"

  # use cloud disk as data disk
  root_password = "${var.instance_password}"

  # we will put all the instances into same vpc and subnet,
  # so they can communicate with each other.
  vpc_id = "${ucloud_vpc.default.id}"

  subnet_id = "${ucloud_subnet.default.id}"

  # this ecurity group to allow HTTP and HTTPS access
  security_group = "${ucloud_security_group.default.id}"

  # easily scale out by count
  count = "${var.instance_count}"
}
