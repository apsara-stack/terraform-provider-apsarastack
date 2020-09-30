---
subcategory: "ECS"
layout: "ApsaraStack"
page_title: "ApsaraStack: apsarastack_cloud_key_pair"
sidebar_current: "docs-apsarastack-resource-key-pair"
description: |-
  Provides a ApsaraStack key pair resource.
---

# apsarastack\_key\_pair

Provides a key pair resource.

## Example Usage

Basic Usage

```
resource "apsarastack_key_pair" "basic" {
  key_name = "terraform-test-key-pair"
}

// Using name prefix to build key pair
resource "apsarastack_key_pair" "prefix" {
  key_name_prefix = "terraform-test-key-pair-prefix"
}

// Import an existing public key to build a apsarastack key pair
resource "apsarastack_key_pair" "publickey" {
  key_name   = "my_public_key"
  public_key = "ssh-rsa AAAAB3Nza12345678qwertyuudsfsg"
}
```
## Argument Reference

The following arguments are supported:

* `key_name` - (ForceNew) The key pair's name.The name must be unique.
* `key_name_prefix` - (ForceNew) The key pair name's prefix. It is conflict with `key_name`. If it is specified, terraform will using it to build the only key name.
* `public_key` - (ForceNew) You can import an existing public key and using ApsaraStack key pair to manage it.
* `key_file` - (ForceNew) The name of file to save your new key pair's private key. Strongly suggest you to specified it when you creating key pair, otherwise, you wouldn't get its private key ever.
* `tags` - (Optional, Available in v1.66.0+) A mapping of tags to assign to the resource.
-> **NOTE:** If `key_name` and `key_name_prefix` are not set, terraform will produce a specified ID to replace.

## Attributes Reference

* `key_name` - The name of the key pair.
* `fingerprint` The finger print of the key pair.
