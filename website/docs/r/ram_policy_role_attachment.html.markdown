---
subcategory: "RAM"
layout: "apsarastack"
page_title: "Apsarastack: apsarastack_ram_policy_role_attachment"
sidebar_current: "docs-apsarastack-resource-ram-policy-role-attachment"
description: |-
Provides a Ram policy role attachment resource.
---

# apsarastack\_ram_policy_role_attachment

Provides a Ram policy role attachment resource.

## Example Usage

```
resource "apsarastack_ram_role" "ram_role" {
  name="testfoorole-001"
  role_policy_document="{\"Version\":\"1\",\"Statement\":[{\"Action\":\"sts:AssumeRole\",\"Effect\":\"Allow\",\"Principal\":{\"Service\":[\"log.aliyuncs.com\"]}}]}"
  description="testing ram role"
}

resource "apsarastack_ram_policy" "ram_policy" {
  name="testfoopolicy-001"
  policy_document="{\"Version\":\"1\",\"Statement\":[{\"Condition\":{},\"Action\":[\"oss:PutObject\"],\"Resource\":[\"acs:oss:*:*:alibaba-cloud/log-collector\",\"acs:oss:*:*:alibaba-cloud/log-collector/*\"],\"Effect\":\"Allow\"}]}"
  description="testing policy"
}
resource "apsarastack_ram_policy_role_attachment" "attach_policy_role" {
  role_name=apsarastack_ram_role.ram_role.name
  policy_name=apsarastack_ram_policy.ram_policy.name
  policy_type="Custom"
}

```
## Argument Reference

The following arguments are supported:

* `role_name` - The name of the role.
* `policy_name` - The name of the policy.
* `policy_type` - The type of the policy.
