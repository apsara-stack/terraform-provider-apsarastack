---
subcategory: "Redis And Memcache (KVStore)"
layout: "apsarastack"
page_title: "Apsarastack: apsarastack_kvstore_connection"
sidebar_current: "docs-apsarastack-resource-kvstore-connection"
description: |-
  Operate the public network ip of the specified resource.
---

# apsarastack\_kvstore\_connection

Operate the public network ip of the specified resource. How to use it, see [What is Resource Alicloud KVStore Connection](https://help.aliyun.com/apsara/enterprise/v_3_16_0_20220117/gpdb/enterprise-ascm-developer-guide/AllocateInstancePublicConnection-1.html?spm=a2c4g.14484438.10001.145).

-> **NOTE:** Available in v1.101.0+.

## Example Usage

Basic Usage

```terraform

variable "name" {
    default = "tf-testAccCheckApsaraStackRKVInstancesDataSource0"
}

data "apsarastack_kvstore_zones" "default"{
	instance_charge_type = "PostPaid"
}

resource "apsarastack_vpc" "default" {
    name       = var.name
    cidr_block = "172.16.0.0/16"
}
resource "apsarastack_vswitch" "default" {
    vpc_id            = apsarastack_vpc.default.id
    cidr_block        = "172.16.0.0/24"
    availability_zone = "${data.apsarastack_kvstore_zones.default.zones.0.id}"
    name              = var.name
}
	
resource "apsarastack_kvstore_instance" "default" {
	instance_name = "tf-testAccKvstoreConnection9302983"
  	instance_class = "redis.master.stand.default"
    vswitch_id     = apsarastack_vswitch.default.id
    private_ip     = "172.16.0.10"
    security_ips   = ["10.0.0.1"]
    instance_type  = "Redis"
    engine_version = "4.0"
}

resource "apsarastack_kvstore_connection" "default" {
  connection_string_prefix = "allocatetestupdate"
  instance_id = "${apsarastack_kvstore_instance.default.id}"
  port = "6371"
}
```

## Argument Reference

The following arguments are supported:
* `connection_string_prefix` - (Required) The prefix of the public endpoint. The prefix can be 8 to 64 characters in length, and can contain lowercase letters and digits. It must start with a lowercase letter.
* `instance_id`- (Required) The ID of the instance.
* `port` - (Required) The service port number of the instance.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of KVStore DBInstance.
* `connection_string` - The public connection string of KVStore DBInstance.

### Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 2 mins) Used when creating the KVStore connection (until it reaches the initial `Normal` status). 
* `update` - (Defaults to 2 mins) Used when updating the KVStore connection (until it reaches the initial `Normal` status). 
* `delete` - (Defaults to 2 mins) Used when deleting the KVStore connection (until it reaches the initial `Normal` status). 

## Import

KVStore connection can be imported using the id, e.g.

```
$ terraform import alicloud_kvstore_connection.example r-abc12345678
```

