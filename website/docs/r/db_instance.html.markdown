---
subcategory: "RDS"
layout: "apsarastack"
page_title: "Apsarastack: apsarastack_db_instance"
sidebar_current: "docs-apsarastack-resource-db-instance"
description: |-
  Provides an RDS instance resource.
---

# apsarastack\_db\_instance

Provides an RDS instance resource. A DB instance is an isolated database
environment in the cloud. A DB instance can contain multiple user-created
databases.

## Example Usage

### Create a RDS MySQL instance

```
variable "name" {
  default = "dbInstanceconfig"
}
variable "creation" {
  default = "Rds"
}
data "apsarastack_zones" "default" {
  available_resource_creation = "${var.creation}"
}
resource "apsarastack_vpc" "default" {
  name       = "${var.name}"
  cidr_block = "172.16.0.0/16"
}
resource "apsarastack_vswitch" "default" {
  vpc_id            = "${apsarastack_vpc.default.id}"
  cidr_block        = "172.16.0.0/24"
  availability_zone = "${data.apsarastack_zones.default.zones.0.id}"
  name              = "${var.name}"
}
resource "apsarastack_db_instance" "default" {
  engine               = "MySQL"
  engine_version       = "5.6"
  instance_type        = "rds.mysql.s2.large"
  instance_storage     = "30"
  instance_name        = "${var.name}"
  vswitch_id           = "${apsarastack_vswitch.default.id}"
  monitoring_period    = "60"
}
```

### Create a RDS MySQL instance with specific parameters

```
resource "apsarastack_vpc" "default" {
  name       = "vpc-123456"
  cidr_block = "172.16.0.0/16"
}

resource "apsarastack_vswitch" "default" {
  vpc_id            = "${apsarastack_vpc.default.id}"
  cidr_block        = "172.16.0.0/24"
  availability_zone = "${data.apsarastack_zones.default.zones.0.id}"
  name              = "vpc-123456"
}

resource "apsarastack_db_instance" "default" {
  engine              = "MySQL"
  engine_version      = "5.6"
  db_instance_class   = "rds.mysql.t1.small"
  db_instance_storage = "10"
  vswitch_id          = "${apsarastack_vswitch.default.id}"
}

resource "apsarastack_db_instance" "default" {
  engine              = "MySQL"
  engine_version      = "5.6"
  db_instance_class   = "rds.mysql.t1.small"
  db_instance_storage = "10"
  parameters {
    name  = "innodb_large_prefix"
    value = "ON"
  }
  parameters {
    name  = "connect_timeout"
    value = "50"
  }
}
```

## Argument Reference

The following arguments are supported:

* `engine` - (Required,ForceNew) Database type. Value options: MySQL, SQLServer, PostgreSQL, and PPAS.
* `engine_version` - (Required,ForceNew) Database version. Value options can refer to the latest docs [CreateDBInstance](https://www.alibabacloud.com/help/doc-detail/26228.htm) `EngineVersion`.
* `instance_type` - (Required) DB Instance type. For details, see [Instance type table](https://www.alibabacloud.com/help/doc-detail/26312.htm).
* `instance_storage` - (Required) User-defined DB instance storage space. Value range:
    - [5, 2000] for MySQL/PostgreSQL/PPAS HA dual node edition;
    - [20,1000] for MySQL 5.7 basic single node edition;
    - [10, 2000] for SQL Server 2008R2;
    - [20,2000] for SQL Server 2012 basic single node edition
    Increase progressively at a rate of 5 GB. For details, see [Instance type table](https://www.alibabacloud.com/help/doc-detail/26312.htm).
    Note: There is extra 5 GB storage for SQL Server Instance and it is not in specified `instance_storage`.

* `instance_name` - (Optional) The name of DB instance. It a string of 2 to 256 characters.
* `zone_id` - (ForceNew) The Zone to launch the DB instance.
If it is a multi-zone and `vswitch_id` is specified, the vswitch must in the one of them.
The multiple zone ID can be retrieved by setting `multi` to "true" in the data source `apsarastack_zones`.
* `vswitch_id` - (ForceNew) The virtual switch ID to launch DB instances in one VPC.
* `security_ips` - (Optional) List of IP addresses allowed to access all databases of an instance. The list contains up to 1,000 IP addresses, separated by commas. Supported formats include 0.0.0.0/0, 10.23.12.24 (IP), and 10.23.12.24/24 (Classless Inter-Domain Routing (CIDR) mode. /24 represents the length of the prefix in an IP address. The range of the prefix length is [1,32]).
* `security_ip_mode` - (Optional)  Valid values are `normal`, `safety`, Default to `normal`. support `safety` switch to high security access mode 
* `parameters` - (Optional) Set of parameters needs to be set after DB instance was launched. Available parameters can refer to the latest docs [View database parameter templates](https://www.alibabacloud.com/help/doc-detail/26284.htm) .
* `force_restart` - (Optional) Set it to true to make some parameter efficient when modifying them. Default to false.
* `tags` - (Optional) A mapping of tags to assign to the resource.
    - Key: It can be up to 64 characters in length. It cannot begin with "aliyun", "acs:", "http://", or "https://". It cannot be a null string.
    - Value: It can be up to 128 characters in length. It cannot begin with "aliyun", "acs:", "http://", or "https://". It can be a null string.

   
* `maintain_time` - (Optional) Maintainable time period format of the instance: HH:MMZ-HH:MMZ (UTC time)

-> **NOTE:** Because of data backup and migration, change DB instance type and storage would cost 15~20 minutes. Please make full preparation before changing them.

## Attributes Reference

The following attributes are exported:

* `id` - The RDS instance ID.
* `port` - RDS database connection port.
* `connection_string` - RDS database connection string.

### Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 20 mins) Used when creating the db instance (until it reaches the initial `Running` status). 
* `update` - (Defaults to 30 mins) Used when updating the db instance (until it reaches the initial `Running` status). 
* `delete` - (Defaults to 20 mins) Used when terminating the db instance. 
