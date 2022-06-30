---
subcategory: "API Gateway"
layout: "apsarastack"
page_title: "Alicloud: apsarastack_api_gateway_vpc_access"
sidebar_current: "docs-apsarastack-resource-api-gateway-vpc-access"
description: |- Provides a Alicloud Api Gateway vpc authorization Resource.
---

# alicloud_api_gateway_app

Provides an vpc authorization resource.This authorizes the API gateway to access your VPC instances.

For information about Api Gateway vpc and how to use it,
see [Set Vpc Access](https://help.aliyun.com/document_detail/400343.html?spm=5176.10695662.1996646101.searchclickresult.67be328fV80qXE)

-> **NOTE:** Terraform will auto build vpc authorization while it uses `apsarastack_api_gateway_vpc_access` to build
vpc.

## Example Usage

Basic Usage

```

variable "name" {
	default = "tf-testAcccn-qingdao-env17-d01ApiGatewayVpcAccess-4857238"
}
	
data "apsarastack_zones" "default" {
	available_disk_category = "cloud_efficiency"
	available_resource_creation= "VSwitch"
}

data "apsarastack_instance_types" "default" {
	availability_zone = "${data.apsarastack_zones.default.zones.0.id}"
}

data "apsarastack_images" "default" {
	name_regex = "^ubuntu"
	most_recent = true
	owners = "system"
}

resource "apsarastack_vpc" "default" {
	vpc_name = "${var.name}"
	cidr_block = "172.16.0.0/12"
}

resource "apsarastack_vswitch" "default" {
	vpc_id = "${apsarastack_vpc.default.id}"
	cidr_block = "172.16.0.0/21"
	availability_zone = "${data.apsarastack_zones.default.zones.0.id}"
}

resource "apsarastack_security_group" "default" {
	name = "${var.name}"
	description = "foo"
	vpc_id = "${apsarastack_vpc.default.id}"
}

resource "apsarastack_instance" "default" {
	vswitch_id = "${apsarastack_vswitch.default.id}"
	image_id = "${data.apsarastack_images.default.images.0.id}"

	# series III
	instance_type = "${data.apsarastack_instance_types.default.instance_types.0.id}"
	system_disk_category = "cloud_efficiency"

	internet_max_bandwidth_out = 5
	security_groups = ["${apsarastack_security_group.default.id}"]
	instance_name = "${var.name}"
}
	
resource "apsarastack_api_gateway_vpc_access" "default" {
  name = "${var.name}"
  vpc_id = "${apsarastack_vpc.default.id}"
  instance_id = "${apsarastack_instance.default.id}"
  port = "8080"
}

```

## Argument Reference

The following arguments are supported:

* `name` - (Required，ForceNew) The name of the vpc authorization.
* `vpc_id` - (Required，ForceNew) The vpc id of the vpc authorization.
* `instance_id` - (Required，ForceNew) ID of the instance in VPC (ECS/Server Load Balance).
* `port` - (Required，ForceNew) ID of the port corresponding to the instance.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the vpc authorization of api gateway.

## Import

Api gateway app can be imported using the id, e.g.

```
$ terraform import apsarastack_api_gateway_vpc_access.example "APiGatewayVpc:vpc-aswcj19ajsz:i-ajdjfsdlf:8080"
```
