---
layout: "ucloud"
page_title: "UCloud: ucloud_db_param_group"
sidebar_current: "docs-ucloud-resource-db-param-group"
description: |-
  Provides a DB param group resource.
---

# ucloud_db_param_group

Provides a DB param group resource.

## Example Usage

## Argument Reference

The following arguments are supported:

* `availability_zone` - (Required) Availability zone where db param groups are located. Such as: "cn-bj-01". You may refer to [list of availability zone](https://docs.ucloud.cn/api/summary/regionlist)当通过id来查询备份时，此参数必填

* `name` - (Required) The name of db param group.
* `description` - (Optional) The description of db param group.
* `src_group_id` - (Required) 源参数组id.
* `engine` - (Required) Database type, possible values are: "mysql", "percona", "postgresql".
* `engine_version` - (Required) The database engine version, possible values are: "5.1", "5.5", "5.6", "5.7", "9.4", "9.6", "10.4".
* `region_flag` - (Optional) 当请求没有填写Zone时，如果指定为true，表示只拉取跨可用区的相关配置文件，否则，拉取所有机房的配置文件（包括每个单可用区和跨可用区）
* `parameter_input` - (Optional) parameter input is a nested type. parameter input documented below.

The attribute (`parameter_input`) support the following:
* `key` - (Required) The key of param.
* `value` - (Required) The value of param.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `parameter_output` - db param groups is a nested type. param groups documented below.

The attribute (`parameter_output`) support the following:

* `key` - (Required) The key of param.
* `value` - (Required) The value of param.
