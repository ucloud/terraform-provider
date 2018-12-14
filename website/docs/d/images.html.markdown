---
layout: "ucloud"
page_title: "UCloud: ucloud_images"
sidebar_current: "docs-ucloud-datasource-images"
description: |-
  Provides a list of available image resources in the current region.
---

# ucloud_images

This data source providers a list of available image resources according to their availability zone, image ID and other fields.

## Example Usage

```hcl
data "ucloud_images" "example" {
    availability_zone = "cn-sh2-02"
    image_type = "Base"
    name_regex = "^CentOS 7.[1-2] 64"
}

output "first" {
    value = "${data.ucloud_images.example.images.0.id}"
}
```

## Argument Reference

The following arguments are supported:

* `availability_zone` - (Optional) Availability zone where instances are located. You may refer to [list of availability zone](https://docs.ucloud.cn/api/summary/regionlist)
* `image_id` - (Optional) The ID of image.
* `name_regex` - (Optional) A regex string to filter resulting images by name. (Such as: `"^CentOS 7.[1-2] 64"` means CentOS 7.1 of 64-bit operating system or CentOS 7.2 of 64-bit operating system, "^Ubuntu 16.04 64" means Ubuntu 16.04 of 64-bit operating system).
* `image_type` - (Optional) The type of image. Possible values are: `"Base"` as standard image, `"Business"` as owned by market place, and `"Custom"` as custom-image, all the image types will be retrieved by default.
* `os_type` - (Optional) The type of OS. Possible values are: `"Linux"` and `"Windows"`, all the OS types will be retrieved by default.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `images` - images is a nested type which documented below.
* `total_count` - Total number of image that satisfy the condition.

The attribute (`images`) support the following:

* `create_time` - The time of creation for EIP, formatted in RFC3339 time string.
* `features` - To identify if any particular feature belongs to the instance, the value is `"NetEnhnced"` as I/O enhanced instance for now.
* `description` - The description of image if any.
* `id` - The ID of image.
* `name` - The name of image.
* `size` - The size of image.
* `type` - The type of image. Possible values are: `"Base"` as standard image, `"Business"` as owned by market place , and `"Custom"` as custom-image, all the image types will be retrieved by default.
* `os_name` - The name of OS.
* `os_type` - The type of OS. Possible values are: `"Linux"` and `"Windows"`, all the OS types will be retrieved by default.
* `status` - The status of image. Possible values are `"Available"`, `"Making"` and `"Unavailable"`.
