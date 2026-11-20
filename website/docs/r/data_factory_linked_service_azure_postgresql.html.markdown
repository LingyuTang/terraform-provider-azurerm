---
subcategory: "Data Factory"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_data_factory_linked_service_azure_postgresql"
description: |-
  Manages a Linked Service (connection) between Azure Database for PostgreSQL and Azure Data Factory.
---

# azurerm_data_factory_linked_service_azure_postgresql

Manages a Linked Service (connection) between Azure Database for PostgreSQL and Azure Data Factory.

Authentication uses a username and a password referenced from Azure Key Vault. Managed identity authentication is not supported.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_data_factory" "example" {
  name                = "example"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
}

data "azurerm_key_vault" "example" {
  name                = "existing-key-vault"
  resource_group_name = "existing-resource-group"
}

resource "azurerm_data_factory_linked_service_key_vault" "example" {
  name            = "example-key-vault"
  data_factory_id = azurerm_data_factory.example.id
  key_vault_id    = data.azurerm_key_vault.example.id
}

resource "azurerm_data_factory_linked_service_azure_postgresql" "example" {
  name            = "example"
  data_factory_id = azurerm_data_factory.example.id
  server          = "example.postgres.database.azure.com"
  port            = 5432
  database_name   = "exampledb"
  ssl_mode        = "5"
  username        = "exampleuser"

  key_vault_password {
    linked_service_name = azurerm_data_factory_linked_service_key_vault.example.name
    secret_name         = "postgresql-password"
  }
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) Specifies the name of the Data Factory Azure Database for PostgreSQL Linked Service. Changing this forces a new resource to be created. Must be unique within a data factory. See the [Microsoft documentation](https://docs.microsoft.com/azure/data-factory/naming-rules) for all restrictions.

* `data_factory_id` - (Required) The Data Factory ID in which to associate the Linked Service with. Changing this forces a new resource to be created.

* `database_name` - (Required) The name of the database on the Azure Database for PostgreSQL server instance.

* `key_vault_password` - (Required) A `key_vault_password` block as defined below. References the password stored in an existing Key Vault through a Key Vault Data Factory Linked Service.

* `server` - (Required) The fully qualified host name of the Azure Database for PostgreSQL server instance.

* `username` - (Required) The username used to authenticate against the Azure Database for PostgreSQL server.

* `annotations` - (Optional) List of tags that can be used for describing the Data Factory Azure Database for PostgreSQL Linked Service.

* `description` - (Optional) The description for the Data Factory Azure Database for PostgreSQL Linked Service.

* `integration_runtime_name` - (Optional) The integration runtime reference to associate with the Data Factory Azure Database for PostgreSQL Linked Service.

* `parameters` - (Optional) A map of parameters to associate with the Data Factory Azure Database for PostgreSQL Linked Service.

* `port` - (Optional) The port number used to connect to the Azure Database for PostgreSQL server instance. Defaults to `5432`.

* `ssl_mode` - (Optional) The SSL connection mode used to connect to the Azure Database for PostgreSQL server. Possible values are `0` (Disabled), `1` (Allow), `2` (Prefer), `3` (Require), `4` (VerifyCA), and `5` (VerifyFull).

---

A `key_vault_password` block supports the following:

* `linked_service_name` - (Required) Specifies the name of an existing Key Vault Data Factory Linked Service.

* `secret_name` - (Required) Specifies the secret name in Azure Key Vault that stores the password.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Data Factory Azure Database for PostgreSQL Linked Service.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/configure#define-operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Data Factory Azure Database for PostgreSQL Linked Service.
* `read` - (Defaults to 5 minutes) Used when retrieving the Data Factory Azure Database for PostgreSQL Linked Service.
* `update` - (Defaults to 30 minutes) Used when updating the Data Factory Azure Database for PostgreSQL Linked Service.
* `delete` - (Defaults to 30 minutes) Used when deleting the Data Factory Azure Database for PostgreSQL Linked Service.

## Import

Data Factory Azure Database for PostgreSQL Linked Service's can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_data_factory_linked_service_azure_postgresql.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resourceGroup1/providers/Microsoft.DataFactory/factories/factory1/linkedservices/linkedservice1
```
