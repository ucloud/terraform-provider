---
layout: "ucloud"
page_title: "UCloud: ucloud_udpn_connection"
sidebar_current: "docs-ucloud-resource-udpn-connection"
description: |-
  Provides an UDPN Connection resource.
---

# ucloud_udpn_connection

UDPN (UCloud Dedicated Private Network)，you can use Dedicated Private Network to achieve high-speed, stable, secure, and dedicated communications between different data centers. The most frequent scenario is to create network connection of clusters across regions.

~> **VPC interconnection** The cross-region Dedicated Private Network must be established if the two VPCs located in different regions are expected to be connected.

~> **Note** The addtional packet head will be added and included in the overall length of packet due to the tunneling UDPN adopted. Since the number of the bytes of packet head is fixed, the bigger data packet is, the less usage will be taken for the packet head.

## Example Usage

```hcl
resource "ucloud_udpn_connection" "example" {
    bandwidth = 2
}
```

## Argument Reference

The following arguments are supported:

* `bandwidth` - (Optional) Maximum bandwidth to the elastic public network, measured in Mbps (Mega bit per second). range from 0 - 1000M. The default value is "1".
* `duration` - (Optional) The duration that you will buy the resource, the default value is "1". It is not required when "Dynamic" (pay by hour), the value is "0" when pay by month and the instance will be vaild till the last day of that month.
* `charge_type` - (Optional) Charge type. Possible values are: "Year" as pay by year, "Month" as pay by month, "Dynamic" as pay by hour. The default value is "Month".
* `peer_region` - (Optional) The correspondent region of dedicated connection, please refer to the region and [availability zone list](https://docs.ucloud.cn/api/summary/regionlist) and [UDPN price list](https://docs.ucloud.cn/network/udpn/udpn_price).

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `create_time` - The time of creation for UDPN connection.
* `expire_time` - The expiration time for UDPN connection.
