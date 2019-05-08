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
