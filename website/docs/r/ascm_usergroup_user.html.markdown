---
subcategory: "ASCM"
layout: "apsarastack"
page_title: "Apsarastack: apsarastack_ascm_usergroup_user"
sidebar_current: "docs-apsarastack-resource-ascm-usergroup_user"
description: |-
  Provides a Ascm usergroup_user resource.
---

# apsarastack\_ascm_usergroup_user

Provides a Ascm usergroup_user resource.

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
 name = "Test_binder"
 parent_id = "1"
}

resource "apsarastack_ascm_user_group" "default" {
 group_name =      "%s"
 organization_id = apsarastack_ascm_organization.default.org_id
}

resource "apsarastack_ascm_user" "default" {
 cellphone_number = "13900000000"
 email = "test@gmail.com"
 display_name = "C2C-DELTA"
 organization_id = apsarastack_ascm_organization.default.org_id
 mobile_nation_code = "91"
 login_name = "User_Role_Test%d"
 login_policy_id = 1
}

resource "apsarastack_ascm_usergroup_user" "default" {
  login_names = ["${apsarastack_ascm_user.default.login_name}"]
  user_group_id = apsarastack_ascm_user_group.default.user_group_id
}

output "org" {
  value = apsarastack_ascm_usergroup_user.default.*
}
```
## Argument Reference

The following arguments are supported:

* `user_group_id` - (Required) group name. 
* `login_names` - (Required) List of user login name.
* `ascm_openapi_endpoint` -For  ascm_openapi_endpoint and how to find it, see [find a ascm_openapi_endpoint](https://help.aliyun.com/apsara/enterprise/v_3_17_0_30393230/apsarabase/enterprise-developer-guide/obtain-the-endpoint-of-a-cloud-service.html?spm=a2c4g.14484438.10001.343)


## Attributes Reference

The following attributes are exported:

* `id` - Login Name of the usergroup_user.