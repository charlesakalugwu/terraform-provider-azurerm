---
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_openshift_cluster"
sidebar_current: "docs-azurerm-resource-container-openshift-cluster"
description: |-
  Manages an Azure Red Hat OpenShift (ARO) Cluster
---

# azurerm_openshift_cluster

Manages an Azure Red Hat OpenShift (ARO) Cluster

~> **Note:** All arguments including secrets will be stored in the raw state as plain-text. [Read more about sensitive data in state](/docs/state/sensitive-data.html).

## Example Usage

This example provisions a basic Managed Azure Red Hat OpenShift Cluster. Other examples of the `azurerm_openshift_cluster` resource can be found in the [examples/openshift](https://github.com/terraform-providers/terraform-provider-azurerm/tree/master/examples/openshift) directory of this repository.

```hcl
variable "prefix" {
  description = "A prefix used for all resources in this example"
  default     = "aro-terraform"
}

variable "location" {
  description = "The Azure Region in which all resources in this example should be provisioned"
  default     = "westeurope"
}

variable "environment" {
  description = "The deployment environment e.g. test, stg, prod"
  default     = "test"
}

variable "client_id" {
  description = "The Client ID for the Service Principal to use for this Managed Azure Red Hat OpenShift Cluster"
}

variable "client_secret" {
  description = "The Client Secret for the Service Principal to use for this Managed Azure Red Hat OpenShift Cluster"
}

variable "customer_admin_group_id" {
  description = "The customer admin group id to use for AAD group memberships that will get synced into the Openshift group osa-customer-admins"
}

data "azurerm_client_config" "current" {}

resource "azurerm_resource_group" "example" {
  name     = "${var.prefix}-example"
  location = var.location

  tags {
    Environment        = var.environment
    TerraformManaged   = true
  }
}

resource "azurerm_openshift_cluster" "example" {
  name                = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name

  openshift_version = "v3.11"

  azure_active_directory = {
    tenant_id               = data.azurerm_client_config.current.tenant_id
    client_id               = var.client_id
    secret                  = var.client_secret
    customer_admin_group_id = var.customer_admin_group_id
  }

  master_pool_profile = {
    count       = 3
    vm_size     = "Standard_D4s_v3"
    subnet_cidr = "10.0.0.0/24"
  }

  infra_pool_profile = {
    count       = 3
    vm_size     = "Standard_D4s_v3"
    subnet_cidr = "10.0.0.0/24"
    os_type     = "linux"
  }

  compute_pool_profile = {
    count       = 3
    vm_size     = "Standard_D4s_v3"
    subnet_cidr = "10.0.0.0/24"
    os_type     = "linux"
  }

  network_profile = {
    vnet_cidr = "10.0.0.0/8"
  }

  tags = {
    Environment        = var.environment
    TerraformManaged   = true
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the Managed OpenShift cluster to create. Changing this forces a new resource to be created.

* `location` - (Required) The location where the Managed OpenShift cluster should be created. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) Specifies the Resource Group where the Managed OpenShift cluster should exist. Changing this forces a new resource to be created.

* `purchase_plan` - (Optional) A `purchase_plan` block.

* `properties` - (Required) A `properties` block.

* `tags` - (Optional) A mapping of tags to assign to the resource.

---

A `purchase_plan` block supports the following:

* `name` - (Optional) The purchase plan ID
* `product` - (Optional) Specifies the product of the image from the marketplace.
* `promotion_code` - (Optional) Specifies a promotion code
* `publisher` - (Optional) The publisher ID

---

A `properties` block supports the following:

* `openshift_version` - (Required) Version of OpenShift specified when creating the cluster.
* `cluster_version` - (Optional) Version of ARO cluster created by the resource provider.
* `network_profile` - (Required) Configuration for OpenShift networking.
* `router_profiles` - (Optional) Configuration for OpenShift router(s).
* `master_pool_profile` - (Required) Configuration for OpenShift master VMs.
* `infra_pool_profile` - (Required) Configuration for OpenShift infra VMs.
* `compute_pool_profile` - (Required) Configuration for OpenShift compute VMs.
* `azure_active_directory` - (Optional) An `azure_active_directory` block. Changing this forces a new resource to be created.

---

A `network_profile` block supports the following:

* `vnet_cidr` - (Required) The CIDR with which the OpenShift cluster's Vnet is configured.
* `vnet_id` - (Optional) The ID of the Vnet created for the OpenShift cluster.
* `peer_vnet_id` - (Optional) The ID of a Vnet to which the OpenShift cluster Vnet should be peered.

---

A `router_profiles` block supports the following:

* `vnet_cidr` - (Required) The CIDR with which the OpenShift cluster's Vnet is configured.
* `vnet_id` - (Optional) The ID of the Vnet created for the OpenShift cluster.
* `peer_vnet_id` - (Optional) The ID of a Vnet to which the OpenShift cluster Vnet should be peered.

---

A `master_pool_profile` block supports the following:

* `name` - (Optional) Unique name of the master pool profile in the context of the subscription and resource group. Only `master` is supported.
* `count` - (Optional)  Number of masters VMs to host the OpenShift cluster control plane. The default value is 3.
* `vm_size` - (Required) Size of master VMs. Possible values include: 'StandardD2sV3', 'StandardD4sV3', 'StandardD8sV3', 'StandardD16sV3', 'StandardD32sV3', 'StandardD64sV3', 'StandardDS4V2', 'StandardDS5V2', 'StandardF8sV2', 'StandardF16sV2', 'StandardF32sV2', 'StandardF64sV2', 'StandardF72sV2', 'StandardF8s', 'StandardF16s', 'StandardE4sV3', 'StandardE8sV3', 'StandardE16sV3', 'StandardE20sV3', 'StandardE32sV3', 'StandardE64sV3', 'StandardGS2', 'StandardGS3', 'StandardGS4', 'StandardGS5', 'StandardDS12V2', 'StandardDS13V2', 'StandardDS14V2', 'StandardDS15V2', 'StandardL4s', 'StandardL8s', 'StandardL16s', 'StandardL32s'.
* `subnet_cidr` - (Required) The CIDR for the subnet in which the master VMs will be provisioned.
* `os_type` - (Optional) Specifies the OS type for the master VMs. Defaults to Linux. Possible values include: `Linux`, `Windows`

---

An `infra_pool_profile` block supports the following:

* `name` - (Required) Unique name of the infra pool profile in the context of the subscription and resource group.
* `count` - (Optional) Number of infra VMs to host the OpenShift infrastructure components. The default value is 3.
* `vm_size` - (Required) Size of infra VMs. Possible values include: 'StandardD2sV3', 'StandardD4sV3', 'StandardD8sV3', 'StandardD16sV3', 'StandardD32sV3', 'StandardD64sV3', 'StandardDS4V2', 'StandardDS5V2', 'StandardF8sV2', 'StandardF16sV2', 'StandardF32sV2', 'StandardF64sV2', 'StandardF72sV2', 'StandardF8s', 'StandardF16s', 'StandardE4sV3', 'StandardE8sV3', 'StandardE16sV3', 'StandardE20sV3', 'StandardE32sV3', 'StandardE64sV3', 'StandardGS2', 'StandardGS3', 'StandardGS4', 'StandardGS5', 'StandardDS12V2', 'StandardDS13V2', 'StandardDS14V2', 'StandardDS15V2', 'StandardL4s', 'StandardL8s', 'StandardL16s', 'StandardL32s'.
* `subnet_cidr` - (Required) The CIDR for the subnet in which the infra VMs will be provisioned.
* `os_type` - (Optional) Specifies the OS type for the infra VMs. Defaults to Linux. Possible values include: `Linux`, `Windows`
* `role` - (Required) Define the role of the AgentPoolProfile. Only `infra` is supported.

---

An `compute_pool_profile` block supports the following:

* `name` - (Required) Unique name of the default compute pool profile in the context of the subscription and resource group.
* `count` - (Optional) Number of default compute VMs to host containers. The default value is 3.
* `vm_size` - (Required) Size of default compute VMs. Possible values include: 'StandardD2sV3', 'StandardD4sV3', 'StandardD8sV3', 'StandardD16sV3', 'StandardD32sV3', 'StandardD64sV3', 'StandardDS4V2', 'StandardDS5V2', 'StandardF8sV2', 'StandardF16sV2', 'StandardF32sV2', 'StandardF64sV2', 'StandardF72sV2', 'StandardF8s', 'StandardF16s', 'StandardE4sV3', 'StandardE8sV3', 'StandardE16sV3', 'StandardE20sV3', 'StandardE32sV3', 'StandardE64sV3', 'StandardGS2', 'StandardGS3', 'StandardGS4', 'StandardGS5', 'StandardDS12V2', 'StandardDS13V2', 'StandardDS14V2', 'StandardDS15V2', 'StandardL4s', 'StandardL8s', 'StandardL16s', 'StandardL32s'.
* `subnet_cidr` - (Required) The CIDR for the subnet in which the default compute VMs will be provisioned.
* `os_type` - (Optional) Specifies the OS type for the default compute VMs. Defaults to Linux. Possible values include: `Linux`, `Windows`
* `role` - (Required) Define the role of the AgentPoolProfile. Only `compute` is supported.

---

An `azure_active_directory` block supports the following:

* `client_id` - (Required) The client id of an Azure Active Directory Application. Changing this forces a new resource to be created.
* `client_secret` - (Required) The client secret of an Azure Active Directory Application. Changing this forces a new resource to be created.
* `tenant_id` - (Required) The tenantId used for the Azure Active Directory Application. If this isn't specified the Tenant ID of the current Subscription is used. Changing this forces a new resource to be created.
* `customer_admin_group_id` - (Optional) The groupId to be granted cluster admin role. Group memberships will get synced into the OpenShift group `osa-customer-admins`

---

## Attributes Reference

The following attributes are exported:

* `id` - The managed OpenShift cluster's resource ID.

* `type` - The managed OpenShift cluster's resource type.

* `application_resource_group` - The managed OpenShift cluster's application resource group.

* `managed_resource_group` - The managed OpenShift cluster's managed resource group.

* `provisioning_state` - The current deployment or provisioning state.

* `public_hostname` - The public hostname of the OpenShift API server.

* `fqdn` - The auto-allocated internal FQDN of the OpenShift API server.

* `vnet_id` - The ID of the Vnet created for the OpenShift cluster.

* `default_router_public_subdomain` - The DNS subdomain for the default OpenShift router profile. The OpenShift master is configured with the public subdomain of the "default" router profile.

* `default_router_fqdn` - The auto-allocated internal FQDN for the default OpenShift router profile.

---

## Import

Managed OpenShift Clusters can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_openshift_cluster.cluster1 /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.ContainerService/managedClusters/cluster1
```
