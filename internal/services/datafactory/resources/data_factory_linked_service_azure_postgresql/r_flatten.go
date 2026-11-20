// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package data_factory_linked_service_azure_postgresql

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/datafactory/2018-06-01/linkedservices"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/datafactory/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func (r Resource) flatten(metadata sdk.ResourceMetaData, id *linkedservices.LinkedServiceId, model *linkedservices.LinkedServiceResource) error {
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
	state := LinkedServiceAzurePostgreSQLModel{
		Name:                   id.LinkedServiceName,
		DataFactoryID:          dataFactoryId.ID(),
		Annotations:            helpers.FlattenTypedLinkedServiceAnnotations(linkedService.Annotations),
		Description:            pointer.From(linkedService.Description),
		IntegrationRuntimeName: helpers.FlattenTypedLinkedServiceIntegrationRuntimeName(linkedService.ConnectVia),
		Parameters:             helpers.FlattenTypedLinkedServiceParameters(linkedService.Parameters),
		KeyVaultPassword:       make([]KeyVaultPassword, 0),
	}

	if props.Password != nil {
		password := KeyVaultPassword{
			LinkedServiceName: props.Password.Store.ReferenceName,
		}
		if secretName, ok := props.Password.SecretName.(string); ok {
			password.SecretName = secretName
		}
		state.KeyVaultPassword = append(state.KeyVaultPassword, password)
	}

	if username, ok := pointer.From(props.Username).(string); ok {
		state.Username = username
	}
	if database, ok := pointer.From(props.Database).(string); ok {
		state.DatabaseName = database
	}
	if server, ok := pointer.From(props.Server).(string); ok {
		state.Server = server
	}

	switch port := pointer.From(props.Port).(type) {
	case int:
		state.Port = int64(port)
	case int64:
		state.Port = port
	case float64:
		state.Port = int64(port)
	}

	switch sslMode := pointer.From(props.SslMode).(type) {
	case string:
		state.SSLMode = sslMode
	case int:
		state.SSLMode = fmt.Sprintf("%d", sslMode)
	case int64:
		state.SSLMode = fmt.Sprintf("%d", sslMode)
	case float64:
		state.SSLMode = fmt.Sprintf("%.0f", sslMode)
	}

	if err := metadata.Encode(&state); err != nil {
		return err
	}

	return pluginsdk.SetResourceIdentityData(metadata.ResourceData, id)
}
