---
subcategory: "Container Service (CS)"
layout: "apsarastack"
page_title: "Apsarastack: apsarastack_cs_kubernetes"
sidebar_current: "docs-apsarastack-resource-cs-kubernetes"
description: |-
  Provides a Apsarastack resource to manage container kubernetes cluster.
---

# apsarastack\_cs\_kubernetes

This resource will help you to manage a Kubernetes Cluster in Apsarastack Kubernetes Service. 

-> **NOTE:** Each kubernetes cluster contains 3 master nodes and those number cannot be changed at now.

-> **NOTE:** Creating kubernetes cluster need to install several packages and it will cost about 15 minutes. Please be patient.
## Example Usage
```$xslt

// If there is not specifying vpc_id, the module will launch a new vpc
resource "apsarastack_vpc" "vpc" {
  name = var.name
  cidr_block = var.cidr_block
}

// According to the vswitch cidr blocks to launch several vswitches
resource "apsarastack_vswitch" "vswitches" {
   vpc_id = "${apsarastack_vpc.default.id}"
   name = var.name
   cidr_block = var.cidr_block
   availability_zone = var.availability_zone
}

resource "apsarastack_cs_kubernetes" "k8s" {
  name = var.name
  vswitch_id = apsarastack_vswitch.vswitches.id
  version = var.version
  master_count = var.master_count
  timeout_mins = var.timeout_mins
  master_disk_category = var.master_disk_category
  master_disk_size = var.master_disk_size
  worker_disk_category = var.worker_disk_category
  worker_disk_size = var.worker_disk_size
  delete_protection = var.delete_protection
  worker_data_disk = var.worker_data_disk
  worker_data_disk_category = var.worker_data_disk_category
  worker_data_disk_size = var.worker_data_disk_size
  new_nat_gateway = var.new_nat_gateway
  slb_internet_enabled = var.slb_internet_enabled
  master_instance_type = var.master_instance_type
  worker_instance_type = var.worker_instance_type
  worker_number         = var.worker_number
  enable_ssh            = var.enable_ssh
  password              = var.password
  pod_cidr              = var.pod_cidr
  service_cidr          = var.service_cidr
}
```

## Argument Reference

The following arguments are supported:

#### Global params
* `name` - (Optional) The kubernetes cluster's name. It is unique in one Apsarastack account.
* `version` - (Optional, Available since 1.70.1) Desired Kubernetes version. If you do not specify a value, the latest available version at resource creation is used and no upgrades will occur except you set a higher version number. The value must be configured and increased to upgrade the version when desired. Downgrades are not supported by ACK.
* `password` - (Required, Sensitive) The password of ssh login cluster node. You have to specify one of `password` `key_name` `kms_encrypted_password` fields.
* `kms_encrypted_password` - (Required, Available in 1.57.1+) An KMS encrypts password used to a cs kubernetes. You have to specify one of `password` `key_name` `kms_encrypted_password` fields.
* `kms_encryption_context` - (Optional, MapString, Available in 1.57.1+) An KMS encryption context used to decrypt `kms_encrypted_password` before creating or updating a cs kubernetes with `kms_encrypted_password`. See [Encryption Context](https://www.alibabacloud.com/help/doc-detail/42975.htm). It is valid when `kms_encrypted_password` is set.
* `user_ca` - (Optional, ForceNew) The path of customized CA cert, you can use this CA to sign client certs to connect your cluster.
* `enable_ssh` - (Optional) Enable login to the node through SSH. default: false 
* `install_cloud_monitor` - (Optional) Install cloud monitor agent on ECS. default: true 
* `cpu_policy` - kubelet cpu policy. options: static|none. default: none.
* `proxy_mode` - Proxy mode is option of kube-proxy. options: iptables|ipvs. default: ipvs.
* `image_id` - Custom Image support. Must based on CentOS7 or AliyunLinux2.
* `user_data` - (Optional, Available in 1.81.0+) Windows instances support batch and PowerShell scripts. If your script file is larger than 1 KB, we recommend that you upload the script to Object Storage Service (OSS) and pull it through the internal endpoint of your OSS bucket.
* `exclude_autoscaler_nodes` - (Optional, Available in 1.88.0+) Exclude autoscaler nodes from `worker_nodes`. default: false 
* `node_name_mode` - (Optional, Available in 1.88.0+) Each node name consists of a prefix, an IP substring, and a suffix. For example, if the node IP address is 192.168.0.55, the prefix is aliyun.com, IP substring length is 5, and the suffix is test, the node name will be aliyun.com00055test. 
* `security_group_id` - (Optional, Available in 1.91.0+) The ID of the security group to which the ECS instances in the cluster belong. If it is not specified, a new Security group will be built.
* `is_enterprise_security_group` - (Optional, Available in 1.91.0+) Enable to create advanced security group. default: false. See [Advanced security group](https://www.alibabacloud.com/help/doc-detail/120621.htm).
* `service_account_issuer` - (Optional, ForceNew, Available in 1.92.0+) The issuer of the Service Account token for [Service Account Token Volume Projection](https://www.alibabacloud.com/help/doc-detail/160384.htm), corresponds to the `iss` field in the token payload. Set this to `"kubernetes.default.svc"` to enable the Token Volume Projection feature (requires specifying `api_audiences` as well).
* `api_audiences` - (Optional, ForceNew, Available in 1.92.0+) A list of API audiences for [Service Account Token Volume Projection](https://www.alibabacloud.com/help/doc-detail/160384.htm). Set this to `["kubernetes.default.svc"]` if you want to enable the Token Volume Projection feature (requires specifying `service_account_issuer` as well.

#### Network
* `pod_cidr` - (Required) [Flannel Specific] The CIDR block for the pod network when using Flannel. 
* `pod_vswitch_ids` - (Required) [Terway Specific] The vswitches for the pod network when using Terway.Be careful the `pod_vswitch_ids` can not equal to `worker_vswtich_ids` or `master_vswtich_ids` but must be in same availability zones.
* `new_nat_gateway` - (Optional) Whether to create a new nat gateway while creating kubernetes cluster. Default to true. Then openapi in Alibaba Cloud are not all on intranet, So turn this option on is a good choice.
* `service_cidr` - (Optional) The CIDR block for the service network. It cannot be duplicated with the VPC CIDR and CIDR used by Kubernetes cluster in VPC, cannot be modified after creation.
* `node_cidr_mask` - (Optional) The node cidr block to specific how many pods can run on single node. 24-28 is allowed. 24 means 2^(32-24)-1=255 and the node can run at most 255 pods. default: 24
* `slb_internet_enabled` - (Optional) Whether to create internet load balancer for API Server. Default to true.

If you want to use `Terway` as CNI network plugin, You need to specific the `pod_vswitch_ids` field and addons with `terway-eniip`.    
If you want to use `Flannel` as CNI network plugin, You need to specific the `pod_cidr` field and addons with `flannel`.

#### Master params
* `master_instance_charge_type` - (Optional) Master payment type. `PrePaid` or `PostPaid`, defaults to `PostPaid`.
* `master_period_unit` - (Optional) Master payment period unit. `Month` or `Week`, defaults to `Month`.
* `master_period` - (Optional) Master payment period. When period unit is `Month`, it can be one of { “1”, “2”, “3”, “4”, “5”, “6”, “7”, “8”, “9”, “12”, “24”, “36”,”48”,”60”}.  When period unit is `Week`, it can be one of {“1”, “2”, “3”, “4”}.
* `master_auto_renew` - (Optional) Enable master payment auto-renew, defaults to false.
* `master_auto_renew_period` - (Optional) Master payment auto-renew period. When period unit is `Month`, it can be one of {“1”, “2”, “3”, “6”, “12”}.  When period unit is `Week`, it can be one of {“1”, “2”, “3”}.
* `master_disk_category` - (Optional) The system disk category of master node. Its valid value are `cloud_ssd` and `cloud_efficiency`. Default to `cloud_efficiency`.
* `master_disk_size` - (Optional) The system disk size of master node. Its valid value range [20~500] in GB. Default to 20.

#### Worker params 
* `worker_number` - (Required) The worker node number of the kubernetes cluster. Default to 3. It is limited up to 50 and if you want to enlarge it, please apply white list or contact with us.
* `worker_disk_size` - (Optional) The system disk size of worker node. Its valid value range [20~32768] in GB. Default to 40.

#### Computed params (No need to configure) 
* `kube_config` - (Optional) The path of kube config, like `~/.kube/config`.
* `client_cert` - (Optional) The path of client certificate, like `~/.kube/client-cert.pem`.
* `client_key` - (Optional) The path of client key, like `~/.kube/client-key.pem`.
* `cluster_ca_cert` - (Optional) The path of cluster ca certificate, like `~/.kube/cluster-ca-cert.pem`
* `availability_zone` - (Optional) The Zone where new kubernetes cluster will be located. If it is not be specified, the `vswitch_ids` should be set, its value will be vswitch's zone.
  
### Timeouts
The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 90 mins) Used when creating the kubernetes cluster (until it reaches the initial `running` status). 
* `update` - (Defaults to 60 mins) Used when activating the kubernetes cluster when necessary during update.
* `delete` - (Defaults to 60 mins) Used when terminating the kubernetes cluster. 

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the container cluster.
* `name` - The name of the container cluster.
* `availability_zone` - The ID of availability zone.
* `vpc_id` - The ID of VPC where the current cluster is located.
* `slb_intranet` - The ID of private load balancer where the current cluster master node is located.
* `security_group_id` - The ID of security group where the current cluster worker node is located.
* `nat_gateway_id` - The ID of nat gateway used to launch kubernetes cluster.
* `master_nodes` - List of cluster master nodes. It contains several attributes to `Block Nodes`.
* `worker_nodes` - List of cluster worker nodes. It contains several attributes to `Block Nodes`.
* `connections` - Map of kubernetes cluster connection information. It contains several attributes to `Block Connections`.
* `version` - The Kubernetes server version for the cluster.
* `worker_ram_role_name` - The RamRole Name attached to worker node.

### Block Nodes
* `id` - ID of the node.
* `name` - Node name.
* `private_ip` - The private IP address of node.


### Block Connections
* `api_server_internet` - API Server Internet endpoint.
* `api_server_intranet` - API Server Intranet endpoint.
* `master_public_ip` - Master node SSH IP address.
* `service_domain` - Service Access Domain.

