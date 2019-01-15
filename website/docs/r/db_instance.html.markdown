---
layout: "ucloud"
page_title: "UCloud: ucloud_db_instance"
sidebar_current: "docs-ucloud-resource-db-instance"
description: |-
  Provides a Database instance resource.
---

# ucloud_db_instance

Provides a Database instance resource.

~> **Note** It takes around 5 mins to shut down the basic database instance (normal version) when making upgrade/degrade(incloud the memory of instance_type and instance_storage), please make the necessary arrangements to your business in advance to prevent any loss of data. In addition, please do confirm if any task pending submission before reset your password, since the password reset will take effect immediately.
## Example Usagek

```hcl
# Query availability zone
data "ucloud_zones" "default" {}

# Create parameter group
data "ucloud_db_parameter_groups" "default" {
  availability_zone = "${data.ucloud_zones.default.zones.0.id}"
  multi_az       = "false"
  engine            = "mysql"
  engine_version    = "5.7"
}

# Create database instance
resource "ucloud_db_instance" "master" {
  availability_zone  = "${data.ucloud_zones.default.zones.0.id}"
  name               = "tf-example-db-instance"
  instance_storage   = 20
  instance_type      = "mysql-ha-1"
  engine             = "mysql"
  engine_version     = "5.7"
  password           = "2018_dbInstance"
  parameter_group_id = "${data.ucloud_db_parameter_groups.default.parameter_groups.0.id}"
  tag                = "tf-example"

  # Backup policy
  backup_begin_time = 4
  backup_count      = 6
  backup_date       = "0111110"
  backup_black_list = ["test.%"]
}
```
## Argument Reference

The following arguments are supported:

* `availability_zone` - (Required) Availability zone where database instances are located. Such as: "cn-bj-01". You may refer to [list of availability zone](https://docs.ucloud.cn/api/summary/regionlist)
* `standby_zone` - (Optional) Availability zone where the standby database instance is located for the high availability database instance with multiple zone; The disaster recovery of data center can be activated by switching to the standby database instance for the high availability database instance.
* `password` - (Optional) The password for the database instance which should have 8-30 characters. It must contain at least 3 items of Capital letters, small letter, numbers and special characters. The special characters include <code>`()~!@#$%^&*-+=_|{}\[]:;'<>,.?/</code>. 
* `engine` - (Required) The type of database engine, possible values are: "mysql", "percona", "postgresql".
* `engine_version` - (Required) The database engine version, possible values are: "5.5", "5.6", "5.7", "9.4", "9.6".
    - 5.5/5.6/5.7 for mysql and percona engine (can not create database slave if under 5.5 version);
    - 9.4/9.6 for postgresql engine.
* `name` - (Optional) The name of the database instance, should have 1 - 63 characters and only support chinese, english, numbers, '-', '_', '.'.
* `instance_storage` - (Optional) Specifies the allocated storage size in gigabytes (GB), range from 20 to 3000GB. The volume adjustment must be a multiple of 10 GB. The maximum disk volume for SSD type are：
    - 500GB if the memory chosen is equal or less than 8GB;
    - 1000GB if the memory chosen is from 12 to 24GB;
    - 2000GB if the memory chosen is 32GB;
    - 3000GB if the memory chosen is equal or more than 48GB.
* `parameter_group_id` - (Optional) The ID of database parameter group. Note: The "parameter_group_id" of the multiple zone database instance should be included in the request for the high availability database instance with multiple zone. When it is changed, the database instance will reboot to make the change take effect.
* `instance_type` - (Required) Specifies the type of database instance with format "engine-type-memory", Possible values are:
    - "mysql","percona" and "postgresql" for engine;
    - "basic" as normal version and  "ha" as high availability version for type of database, thereinto, high availability version use the dual main hot standby structure which can thoroughly solved the issue of unavailable database caused by the system downtime or hardware failure, the "ha" version only supports "mysql" and "percona" engine, the standard version only supports the "postgresql" engine.
    - possible values for memory are: 1, 2, 4, 6, 8, 12, 16, 24, 32, 48, 64GB.
* `port` - (Optional) The port on which the database accepts connections, the default port is 3306 for mysql and percona and 5432 for postgresql.
* `instance_charge_type` - (Optional) The charge type of database instance, possible values are: "Year", "Month" and "Dynamic" as pay by hour (specific permission required). the dafault is "Month".
* `instance_duration` - (Optional) The duration that you will buy the resource, the default value is "1". It is not required when "Dynamic" (pay by hour), the value is "0" when pay by month and the instance will be vaild till the last day of that month.
* `vpc_id` - (Optional) The ID of VPC linked to the database instances.
* `subnet_id` - (Optional) The ID of subnet.
* `backup_count` - (Optional) Specifies the number of backup saved per week, it is 7 backups saved per week by default.
* `backup_begin_time` - (Optional) Specifies when the backup starts, measured in hour, it starts at one o'clock of 1, 2, 3, 4 in the morning by default.
* `backup_date` - (Optional) Specifies whether the backup took place from Sunday to Saturday by displaying 7 digits. 0 stands for backup disbaled and 1 stands for backup enabled. The rightmost digit specifies whether the backup took place on Sunday, and the digits from right to left specify whether the backup took place from Monday to Saturday, it's mandatory required to backup twice per week at least. such as: digits "1100000" stands for the backup took place on Saturday and Friday.
* `backup_id` - (Optional) The ID of backup set of database instance, The instance is created based on a backup set if the ID is specified, otherwise the ID is set to "null". Please note that the "availability_zone ","engine" and "engine_version" requested must be identical with the backup set when performing recovery from backup set.
* `backup_black_list` - (Optional) The backup for database such as "test.%" or table such as "city.address" specified in the black lists are not supprted.
* `tag` - (Optional) A tag to assign to the instance. The default value is "Default" (means no tag assigned).

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `status` - Specifies the status of database, possible values are: "Init","Fail", "Starting", "Running", "Shutdown", "shutoff", "Delete", "Upgrading", "Promoting", "Recovering" and "Recover fail".
* `create_time` - The creation time of database, formatted by RFC3339 time string.
* `expire_time` - The expiration time of database, formatted by RFC3339 time string.
* `modify_time` - The modification time of database, formatted by RFC3339 time string.
