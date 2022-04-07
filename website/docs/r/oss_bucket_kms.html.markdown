---
subcategory: "OSS"
layout: "apsarastack"
page_title: "Apsarastack: apsarastack_oss_bucket_kms"
sidebar_current: "docs-apsarastack-resource-oss-bucket-kms"
description: |-
  Provides a resource to create an oss bucket kms.
---

# apsarastack\_oss\_bucket

Provides a resource to create an oss bucket and set its attribution.

-> **NOTE:** The bucket namespace is shared by all users of the OSS system. Please set bucket name as unique as possible.


## Example Usage

Private Bucket

```
resource "apsarastack_oss_bucket" "default" {
  bucket = "sample_bucket"
  acl    = "public-read"
}

resource "apsarastack_oss_bucket_kms" "default" {
  bucket = "${apsarastack_oss_bucket.default.bucket}"
  sse_algorithm    = "KMS"
  kms_data_encryption = "SM4"
}
```

## Argument Reference

The following arguments are supported:

* `sse_algorithm` - (require) It Can be "KMS", and "AES256".
* `kms_data_encryption` - (require) When the encryption method is SSE-KMS, OSS uses the default encryption algorithm AES256 for encryption. If you want to encrypt with encryption algorithm SM4, please specify through this option.

## Attributes Reference

The following attributes are exported:

* `sse_algorithm` - (require) It Can be "KMS", and "AES256".
* `kms_data_encryption` - (require) When the encryption method is SSE-KMS, OSS uses the default encryption algorithm AES256 for encryption. If you want to encrypt with encryption algorithm SM4, please specify through this option.


