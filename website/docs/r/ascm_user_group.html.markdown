---
subcategory: "ASCM"
layout: "apsarastack"
page_title: "Apsarastack: apsarastack_ascm_user_group"
sidebar_current: "docs-apsarastack-resource-ascm-user_group"
description: |-
  Provides a Ascm user group resource.
---

# apsarastack\_ascm_user_group

Provides a Ascm user group resource.

## Example Usage
provider "apsarastack" {
domain    = "***"
access_key  = "***"
secret_key  = "***"
region   = "***"
proxy   = "***"
protocol   = "***"
insecure   = "***"
resource_group_set_name = "***"
ascm_openapi_endpoint = "***"
}
```
resource "apsarastack_ascm_organization" "default" {
  name = "Dummy_Test_1"
}

resource "apsarastack_ascm_user_group" "default" {
   group_name = "test"
   organization_id = apsarastack_ascm_organization.default.org_id
   role_in_ids =   []string{"2", "6"}
}

output "org" {
  value = apsarastack_ascm_user_group.default.*
}
```
## Argument Reference

The following arguments are supported:

* `group_name` - (Required) group name. 
* `organization_id` - (Required) User Organization ID.
* `role_in_ids` - (Optional) ascm role id.
* `ascm_openapi_endpoint` -For  ascm_openapi_endpoint and how to find it, see [find a ascm_openapi_endpoint](https://help.aliyun.com/apsara/enterprise/v_3_17_0_30393230/apsarabase/enterprise-developer-guide/obtain-the-endpoint-of-a-cloud-service.html?spm=a2c4g.14484438.10001.343)

## Attributes Reference

The following attributes are exported:

* `id` - Login Name of the user group.