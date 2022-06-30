---
subcategory: "RDS"
layout: "apsarastack"
page_title: "Apsarastack: apsarastack_rds_parameter_group"
sidebar_current: "docs-apsarastack-resource-rds-parameter-group"
description: |-
  Provides a Apsarastack RDS Parameter Group resource.
---

# apsarastack\_rds\_parameter\_group

Provides a RDS Parameter Group resource.

For information about RDS Parameter Group and how to use it, see [What is Parameter Group](https://help.aliyun.com/document_detail/207419.html?spm=5176.21213303.J_6704733920.7.34ac53c9jZzAiI&scm=20140722.S_help%40%40%E6%96%87%E6%A1%A3%40%40207419._.ID_help%40%40%E6%96%87%E6%A1%A3%40%40207419-RL_CreateParameterGroup-LOC_main-OR_ser-V_2-P0_0).

-> **NOTE:** Available in v1.119.0+.

## Example Usage

Basic Usage

```terraform

variable "name" {
			default = "tf_testAccApsaraStackRdsParameterGroup12586"
		}
resource "apsarastack_rds_parameter_group" "default" {
  engine_version = "5.7"
  param_detail {
    param_name = "back_log"
    param_value = "3000"
  }
  param_detail {
    param_name = "wait_timeout"
    param_value = "86400"
  }
  
  parameter_group_desc = "test"
  parameter_group_name = "${var.name}"
  engine = "mysql"
}
```

## Argument Reference

The following arguments are supported:

* `engine` - (Required, ForceNew) The database engine. Valid values: `mysql`, `mariadb`.
* `engine_version` - (Required, ForceNew) The version of the database engine. Valid values: mysql: `5.1`, `5.5`, `5.6`, `5.7`, `8.0`; mariadb: `10.3`.
* `param_detail` - (Required) Parameter list.
* `parameter_group_desc` - (Optional) The description of the parameter template.
* `parameter_group_name` - (Required) The name of the parameter template.

#### Block parameter_detail

The param_detail supports the following: 

* `param_name` - (Required) The name of a parameter.
* `param_value` - (Required) The value of a parameter.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of Parameter Group.

## Import

RDS Parameter Group can be imported using the id, e.g.

```
$ terraform import alicloud_rds_parameter_group.example <id>
```
