---
layout: "ucloud"
page_title: "Provider: UCloud"
sidebar_current: "docs-ucloud-index"
description: |-
  The UCloud provider is used to interact with many resources supported by UCloud. The provider needs to be configured with the proper credentials before it can be used.
---

# UCloud Provider

~> **NOTE:** This guide requires an avaliable UCloud account or sub-account with project to create resources.

The UCloud provider is used to interact with the
resources supported by UCloud. The provider needs to be configured
with the proper credentials before it can be used.

Use the navigation to the left to read about the available resources.

## Example Usage

```hcl
# Configure the UCloud Provider
provider "ucloud" {
  public_key = "${var.ucloud_public_key}"
  private_key = "${var.ucloud_private_key}"
  project_id = "${var.ucloud_project_id}"
  region     = "cn-sh2"
}

# Query availability zone
data "ucloud_zones" "default" {
}

# Query image
data "ucloud_images" "default" {
    availability_zone = "${data.ucloud_zones.default.zones.0.id}"
    os_type = "Linux"
}

# Create security group
resource "ucloud_security_group" "default" {
    name = "tf-example-eip"
    tag  = "tf-example"

    rules {
        port_range = "80"
        protocol   = "TCP"
        cidr_block = "192.168.0.0/16"
        policy     = "ACCEPT"
    }
}

# Create a web server
resource "ucloud_instance" "web" {
    instance_type     = "n-standard-1"
    availability_zone = "${data.ucloud_zones.default.zones.0.id}"
    image_id = "${data.ucloud_images.default.images.0.id}"

    root_password      = "wA1234567"
    security_group     = "${ucloud_security_group.default.id}"

    name              = "tf-example-eip"
    tag               = "tf-example"
}
```

## Authentication

The UCloud provider offers a flexible means of providing credentials for
authentication. The following methods are supported, in this order, and
explained below:

- Static credentials
- Environment variables

### Static credentials

Static credentials can be provided by adding an `public_key` and `private_key` in-line in the
UCloud provider block:

Usage:

```hcl
provider "ucloud" {
  public_key = "your_public_key"
  private_key = "your_private_key"
  project_id = "your_project_id"
  region     = "cn-sh2"
}
```

### Environment variables

You can provide your credentials via `UCLOUD_PUBLIC_KEY` and `UCLOUD_PRIVATE_KEY`
environment variables, representing your UCloud public key and private key respectively.
`UCLOUD_REGION` and `UCLOUD_PROJECT_ID` are also used, if applicable:

```hcl
provider "ucloud" {}
```

Usage:

```hcl
$ export UCLOUD_PUBLIC_KEY="your_public_key"
$ export UCLOUD_PRIVATE_KEY="your_private_key"
$ export UCLOUD_REGION="cn-sh2"
$ export UCLOUD_PROJECT_ID="org-xxx"

$ terraform plan
```

## Argument Reference

In addition to [generic `provider` arguments](https://www.terraform.io/docs/configuration/providers.html)
(e.g. `alias` and `version`), the following arguments are supported in the UCloud
 `provider` block:

* `public_key` - (Required) This is the UCloud public key. It must be provided, but
  it can also be sourced from the `UCLOUD_PUBLIC_KEY` environment variable.

* `private_key` - (Required) This is the UCloud private key. It must be provided, but
  it can also be sourced from the `UCLOUD_PRIVATE_KEY` environment variable.

* `region` - (Required) This is the UCloud region. It must be provided, but
  it can also be sourced from the `UCLOUD_REGION` environment variables.

* `project_id` - (Required) This is the UCloud project id. It must be provided, but
  it can also be sourced from the `UCLOUD_PROJECT_ID` environment variables.

* `max_retries` - (Optional) This is the max retry attempts number. Default max retry attempts number is `"0"`.

* `insecure` - (Optional) This is a switch to disable/enable https. (Default: `"false"`, means enable https).

## Testing

Credentials must be provided via the `UCLOUD_PUBLIC_KEY`, `UCLOUD_PRIVATE_KEY`, `UCLOUD_PROJECT_ID` environment variables in order to run acceptance tests.
