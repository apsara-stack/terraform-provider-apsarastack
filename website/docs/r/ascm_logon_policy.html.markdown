---
subcategory: "ASCM"
layout: "apsarastack"
page_title: "Apsarastack: apsarastack_ascm_logon_policy"
sidebar_current: "docs-apsarastack_ascm_logon_policy"
description: |-
  Provides a Apsarastack Logon Policy resource.
---
# apsarastack\_ascm_logon_policy

Provides a Apsarastack Logon Policy resource.

Basic Usage

```
resource "apsarastack_ascm_logon_policy" "login" {
  name="test_foo"
  description="testing purpose"
  rule="ALLOW"
}
```
## Argument Reference

The following arguments are supported:
* `name` - (Required) The name of the Logon Policy. Defaults to null.
* `description` - (Optional) The Logon Policy description. Defaults to null.
* `rule` - (Optional) The Rule for the Logon Policy.

## Attributes Reference

The following attributes are exported:

* `name` - The name of the Logon Policy. Defaults to null.
* `description` - (Optional) The Logon Policy description. Defaults to null.
* `rule` - (Optional) The Rule for the Logon Policy.
* `policy_id` - The ID of the logon policy created.

