// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package data_factory_linked_service_azure_postgresql

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/datafactory/2018-06-01/linkedservices"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

const azurermDataFactoryLinkedServiceAzurePostgreSQLResourceName = "azurerm_data_factory_linked_service_azure_postgresql"

type Resource struct{}

var (
	_ sdk.ResourceWithUpdate   = Resource{}
	_ sdk.ResourceWithIdentity = Resource{}
)

func (r Resource) ResourceType() string {
	return azurermDataFactoryLinkedServiceAzurePostgreSQLResourceName
}

func (r Resource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return linkedservices.ValidateLinkedServiceID
}

func (r Resource) Identity() resourceids.ResourceId {
	return &linkedservices.LinkedServiceId{}
}
