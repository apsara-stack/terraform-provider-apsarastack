---
subcategory: "Container Registry (CR)"
layout: "apsarastack"
page_title: "Apsarastack: apsarastack_cr_ee_repo"
sidebar_current: "docs-apsarastack-resource-cr-ee-repo"
description: |-
  Provides a Apsarastack resource to manage Container Registry Enterprise Edition repositories.
---

# apsarastack\_cr\_ee\_repo

This resource will help you to manager Container Registry Enterprise Edition repositories.


## Example Usage

Basic Usage

```
resource "apsarastack_cr_ee_namespace" "my-namespace" {
  instance_id        = "cri-xxx"
  name               = "my-namespace"
  auto_create        = false
  default_visibility = "PUBLIC"
}

resource "apsarastack_cr_ee_repo" "my-repo" {
  instance_id = apsarastack_cr_ee_namespace.my-namespace.instance_id
  namespace   = apsarastack_cr_ee_namespace.my-namespace.name
  name        = "my-repo"
  summary     = "this is summary of my new repo"
  repo_type   = "PUBLIC"
  detail      = "this is a public repo"
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required, ForceNew) ID of Container Registry Enterprise Edition instance.
* `namespace` - (Required, ForceNew) Name of Container Registry Enterprise Edition namespace where repository is located. It can contain 2 to 30 characters.
* `name` - (Required, ForceNew) Name of Container Registry Enterprise Edition repository. It can contain 2 to 64 characters.
* `summary` - (Required) The repository general information. It can contain 1 to 100 characters.
* `repo_type` - (Required) `PUBLIC` or `PRIVATE`, repo's visibility.
* `detail` - (Optional) The repository specific information. MarkDown format is supported, and the length limit is 2000.

## Attributes Reference

The following attributes are exported:

* `id` - The resource id of Container Registry Enterprise Edition repository. The value is in format `{instance_id}:{namespace}:{repository}`.
* `repo_id` - The uuid of Container Registry Enterprise Edition repository.

