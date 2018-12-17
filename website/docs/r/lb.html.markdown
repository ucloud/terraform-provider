---
layout: "ucloud"
page_title: "UCloud: ucloud_lb"
sidebar_current: "docs-ucloud-resource-lb"
description: |-
  Provides a Load Balancer resource.
---

# ucloud_lb

Provides a Load Balancer resource.

## Example Usage

```hcl
resource "ucloud_lb" "web" {
    name = "tf-example-lb"
    tag  = "tf-example"
}
```

## Argument Reference

The following arguments are supported:

* `internal` - (Optional) Indicate whether the load balancer is intranet.
* `charge_type` - (Optional) Charge type of load balancer. Possible values are: `"year"` as pay by year, `"month"` as pay by month, `"dynamic"` as pay by hour (specific permission required). (Default: `"month"`).
* `name` - (Optional) The name of the load balancer. If not specified, terraform will autogenerate a name beginning with `"tf-lb"`.
* `vpc_id` - (Optional) The ID of the VPC linked to the Load Balancers, This argumnet is not required if default VPC.
* `subnet_id` - (Optional) The ID of subnet that intrant load balancer belongs to. This argumnet is not required if default subnet.
* `tag` - (Optional) A mapping of tags to assign to the load balancer, which contains 1-63 characters and only support Chinese, English, numbers, '-', '_', '.'. (Default: `"Default"`).
* `remark` - (Optional) The remarks of the load balancer. (Default: is `""`).

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `create_time` - The time of creation for load balancer, formatted in RFC3339 time string.
* `expire_time` - The expiration time for load balancer, formatted in RFC3339 time string.
* `ip_set` - It is a nested type which documented below.
* `private_ip` - The IP address of intranet IP. It is `""` if `internal` is `false`.

The attribute (`ip_set`) support the following:

* `internet_type` - Type of Elastic IP routes.
* `ip` - Elastic IP address.
