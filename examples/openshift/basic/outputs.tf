output "id" {
  value = "${azurerm_openshift_cluster.example.id}"
}

output "type" {
  value = "${azurerm_openshift_cluster.example.type}"
}

output "provisioning_state" {
  value = "${azurerm_openshift_cluster.example.properties.provisioning_state}"
}

output "fqdn" {
  value = "${azurerm_openshift_cluster.example.properties.fqdn}"
}

output "public_hostname" {
  value = "${azurerm_openshift_cluster.example.properties.public_hostname}"
}

output "vnet_id" {
  value = "${azurerm_openshift_cluster.example.properties.network_profile.vnet_id}"
}

output "default_router_profile_fqdn" {
  value = "${azurerm_openshift_cluster.example.properties.router_profiles.default.fqdn}"
}

output "default_router_profile_public_subdomain" {
  value = "${azurerm_openshift_cluster.example.properties.router_profiles.default.public_subdomain}"
}
