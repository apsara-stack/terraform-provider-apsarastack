---
subcategory: "Log Service (SLS)"
layout: "apsarastack"
page_title: "Apsarastack: apsarastack_log_project"
sidebar_current: "docs-apsarastack-resource-log-project"
description: |-
  Provides a Apsarastack log project resource.
---

# apsarastack\_log\_project

The project is the resource management unit in Log Service and is used to isolate and control resources.
You can manage all the logs and the related log sources of an application by using projects.

## Example Usage

Basic Usage

```
resource "apsarastack_log_project" "example" {
  name        = "tf-log"
  description = "created by terraform"
}
```


## Argument Reference

The following arguments are supported:

* `name` - (Required, ForceNew) The name of the log project. It is the only in one Apsarastack account.
* `description` - (Optional) Description of the log project.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the log project. It sames as its name.
* `name` - Log project name.
* `description` - Log project description.


