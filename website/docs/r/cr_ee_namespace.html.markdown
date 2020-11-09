---
subcategory: "Container Registry(CR)"
layout: "apsarastack"
page_title: "Apsarastack: apsarastack_cr_ee_namespace"
sidebar_current: "docs-apsarastack-resource-cr-ee-namespace"
description: |-
  Provides a Apsarastack resource to manage Container Registry Enterprise Edition namespaces.
---

# apsarastack\_cr\_ee\_namespace

This resource will help you to manager Container Registry Enterprise Edition namespaces.

## Example Usage

Basic Usage

```
resource "apsarastack_cr_ee_namespace" "my-namespace" {
  instance_id        = "cri-xxx"
  name               = "my-namespace"
  auto_create        = false
  default_visibility = "PUBLIC"
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required, ForceNew) ID of Container Registry Enterprise Edition instance.
* `name` - (Required, ForceNew) Name of Container Registry Enterprise Edition namespace. It can contain 2 to 30 characters.
* `auto_create` - (Required) Boolean, when it set to true, repositories are automatically created when pushing new images. If it set to false, you create repository for images before pushing.
* `default_visibility` - (Required) `PUBLIC` or `PRIVATE`, default repository visibility in this namespace.

## Attributes Reference

The following attributes are exported:

* `id` - ID of Container Registry Enterprise Edition namespace. The value is in format `{instance_id}:{namespace}` .

