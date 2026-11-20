// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package data_factory_linked_service_azure_postgresql

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/datafactory/2018-06-01/linkedservices"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/datafactory/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func resourceDataFactoryLinkedServiceAzurePostgreSQLFlatten(d *pluginsdk.ResourceData, id *linkedservices.LinkedServiceId, model *linkedservices.LinkedServiceResource) error {
	if model == nil {
		return fmt.Errorf("retrieving %s: `model` was nil", id)
	}
	if model.Properties == nil {
		return fmt.Errorf("retrieving %s: `properties` was nil", id)
	}

	linkedService, ok := model.Properties.(linkedservices.AzurePostgreSqlLinkedService)
	if !ok {
		return fmt.Errorf("classifying %s: expected `LinkedServiceAzurePostgreSql` received `%T`", id, model.Properties)
	}

	props := linkedService.TypeProperties
	dataFactoryId := linkedservices.NewFactoryID(id.SubscriptionId, id.ResourceGroupName, id.FactoryName)
	state := map[string]any{
		"name":                     id.LinkedServiceName,
		"data_factory_id":          dataFactoryId.ID(),
		"annotations":              helpers.FlattenTypedLinkedServiceAnnotations(linkedService.Annotations),
		"description":              pointer.From(linkedService.Description),
		"integration_runtime_name": helpers.FlattenTypedLinkedServiceIntegrationRuntimeName(linkedService.ConnectVia),
		"parameters":               helpers.FlattenTypedLinkedServiceParameters(linkedService.Parameters),
		"key_vault_password":       helpers.FlattenTypedLinkedServiceKeyVaultPassword(props.Password),
		"username":                 "",
		"database_name":            "",
		"server":                   "",
		"port":                     0,
		"ssl_mode":                 "",
	}

	if username, ok := pointer.From(props.Username).(string); ok {
		state["username"] = username
	}
	if database, ok := pointer.From(props.Database).(string); ok {
		state["database_name"] = database
	}
	if server, ok := pointer.From(props.Server).(string); ok {
		state["server"] = server
	}

	switch port := pointer.From(props.Port).(type) {
	case int:
		state["port"] = port
	case int64:
		state["port"] = int(port)
	case float64:
		state["port"] = int(port)
	}

	switch sslMode := pointer.From(props.SslMode).(type) {
	case string:
		state["ssl_mode"] = sslMode
	case int:
		state["ssl_mode"] = fmt.Sprintf("%d", sslMode)
	case int64:
		state["ssl_mode"] = fmt.Sprintf("%d", sslMode)
	case float64:
		state["ssl_mode"] = fmt.Sprintf("%.0f", sslMode)
	}

	for name, value := range state {
		if err := d.Set(name, value); err != nil {
			return fmt.Errorf("setting `%s`: %+v", name, err)
		}
	}

	return pluginsdk.SetResourceIdentityData(d, id)
}
