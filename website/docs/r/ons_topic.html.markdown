---
subcategory: "RocketMQ"
layout: "apsarastack"
page_title: "Apsarastack: apsarastack_ons_topic"
sidebar_current: "docs-apsarastack-resource-ons-topic"
description: |-
  Provides a apsarastack ONS Topic resource.
---

# apsarastack\_ons\_topic

Provides an ONS topic resource.


## Example Usage

Basic Usage

```
variable "name" {
  default = "onsInstanceName"
}

variable "topic" {
  default = "onsTopicName"
}

resource "apsarastack_ons_instance" "default" {
  name = "${var.name}"
  remark = "default_ons_instance_remark"
}

resource "apsarastack_ons_topic" "default" {
  topic = "${var.topic}"
  instance_id = "${apsarastack_ons_instance.default.id}"
  message_type = 0
  remark = "dafault_ons_topic_remark"
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required) ID of the ONS Instance that owns the topics.
* `topic` - (Required) Name of the topic. Two topics on a single instance cannot have the same name and the name cannot start with 'GID' or 'CID'. The length cannot exceed 64 characters.
* `message_type` - (Required) The type of the message.
* `remark` - (Optional) This attribute is a concise description of topic. The length cannot exceed 128.
* `perm` - (Optional) This attribute is used to set the read-write mode for the topic.

## Attributes Reference

The following attributes are exported:

* `id` - The `key` of the resource supplied above. The value is formulated as `<instance_id>:<topic>`.

