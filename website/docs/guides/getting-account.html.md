---
layout: "apsarastack"
page_title: "ApsaraStack Cloud Account Guide"
sidebar_current: "docs-apsarastack-guide-apsarastack-account-guide"
description: |-
  Sign up ApsaraStack Cloud and distinguish its type.
---

# Getting ApsaraStack Cloud Account

The ApsaraStack Cloud has three accounts: International-Site Account, China-Site Account and JP-Site Account.
In most case, the three accounts have no different about creating ApsaraStack Cloud resources.
But, based on some internal reason, when using terraform to manage cloud resources,
a few products and resources have some limitations are applied to different accounts.
We will show the limitations gradually to help you avoid some needless errors.

## Sign Up ApsaraStack Cloud Internation-Site Account

-> **Warning:** At present, terraform can not use internation-site to open `Subscription`
resources which instance charge type is "PrePaid"

If you want to sign up a Internation-Site account, you can go to [ApsaraStack Cloud Internationl-Site Website](https://www.alibabacloud.com/)
to finish register. For more account register details, see [Sign up with ApsaraStack Cloud](https://www.alibabacloud.com/help/doc-detail/50482.html)

## Sign Up ApsaraStack Cloud China-Site Account

China-Site has different access website. If you want to sign up a China-Site account, you can go to
[ApsaraStack Cloud China-Site Website](https://www.aliyun.com/) to finish register.
For more account register details, see [Sign up with ApsaraStack Cloud](https://help.aliyun.com/knowledge_detail/37195.html)

## Sign Up ApsaraStack Cloud JP-Site Account

JP-Site(Japan-Site) also has a alone access website. If you want to sign up a China-Site account, you can go to
[ApsaraStack Cloud JP-Site Website](https://jp.alibabacloud.com/) to finish register.
For more account register details, see [Sign up with ApsaraStack Cloud](https://www.alibabacloud.com/help/doc-detail/50482.html)

## How to distinguish my account site type

There is a simple method to distinguish an ApsaraStack Cloud account belongs to Internation-Site, China-Site or JP-Site:
An account can only access the corresponding site, that are Internation-Site account can only login [Internationl-Site Website](https://www.alibabacloud.com/),
China-Site account only login [China-Site Website](https://www.aliyun.com/) and JP-Site can only login [JP-Site Website](https://jp.alibabacloud.com/).

