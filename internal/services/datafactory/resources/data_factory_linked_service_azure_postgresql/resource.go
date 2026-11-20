// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package data_factory_linked_service_azure_postgresql

import (
	"time"

	"github.com/hashicorp/go-azure-sdk/resource-manager/datafactory/2018-06-01/linkedservices"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

const azurermDataFactoryLinkedServiceAzurePostgreSQLResourceName = "azurerm_data_factory_linked_service_azure_postgresql"

func RegisterResource() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceDataFactoryLinkedServiceAzurePostgreSQLCreateUpdate,
		Read:   resourceDataFactoryLinkedServiceAzurePostgreSQLRead,
		Update: resourceDataFactoryLinkedServiceAzurePostgreSQLCreateUpdate,
		Delete: resourceDataFactoryLinkedServiceAzurePostgreSQLDelete,

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Identity: &schema.ResourceIdentity{
			SchemaFunc: pluginsdk.GenerateIdentitySchema(&linkedservices.LinkedServiceId{}),
		},
		Importer: pluginsdk.ImporterValidatingIdentity(&linkedservices.LinkedServiceId{}),
		Schema:   dataFactoryLinkedServiceAzurePostgreSQLSchema(),
	}
}
