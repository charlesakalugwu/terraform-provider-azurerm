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
