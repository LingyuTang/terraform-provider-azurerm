// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package data_factory_linked_service_azure_postgresql

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/datafactory/2018-06-01/linkedservices"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/datafactory/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func (r Resource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func:    r.createUpdate,
	}
}

func (r Resource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func:    r.createUpdate,
	}
}

func (r Resource) createUpdate(ctx context.Context, metadata sdk.ResourceMetaData) error {
	client := metadata.Client.DataFactory.LinkedServicesClient

	var config LinkedServiceAzurePostgreSQLModel
	if err := metadata.Decode(&config); err != nil {
		return fmt.Errorf("decoding: %+v", err)
	}

	dataFactoryId, err := linkedservices.ParseFactoryID(config.DataFactoryID)
	if err != nil {
		return err
	}

	id := linkedservices.NewLinkedServiceID(metadata.Client.Account.SubscriptionId, dataFactoryId.ResourceGroupName, dataFactoryId.FactoryName, config.Name)
	if metadata.ResourceData.IsNewResource() {
		if !metadata.Client.Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
			existing, err := client.Get(ctx, id, linkedservices.DefaultGetOperationOptions())
			if err != nil && !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
			}
			if !response.WasNotFound(existing.HttpResponse) {
				return metadata.ResourceRequiresImport(r.ResourceType(), id)
			}
		}
	} else {
		existingId, err := linkedservices.ParseLinkedServiceID(metadata.ResourceData.Id())
		if err != nil {
			return err
		}
		id = *existingId
	}

	annotations := make([]any, len(config.Annotations))
	for i, annotation := range config.Annotations {
		annotations[i] = annotation
	}

	parameters := linkedservices.LinkedServiceResource{
		Properties: linkedservices.AzurePostgreSqlLinkedService{
			Annotations:    helpers.ExpandTypedLinkedServiceAnnotations(annotations),
			ConnectVia:     helpers.ExpandTypedLinkedServiceIntegrationRuntimeName(config.IntegrationRuntimeName),
			Description:    pointer.To(config.Description),
			Parameters:     helpers.ExpandTypedLinkedServiceParameters(config.Parameters),
			TypeProperties: expandLinkedServiceAzurePostgreSQLTypeProperties(config),
			Version:        pointer.To("2.0"),
		},
	}

	if _, err := client.CreateOrUpdate(ctx, id, parameters, linkedservices.DefaultCreateOrUpdateOperationOptions()); err != nil {
		return fmt.Errorf("creating/updating %s: %+v", id, err)
	}

	metadata.SetID(id)
	return pluginsdk.SetResourceIdentityData(metadata.ResourceData, &id)
}

func expandLinkedServiceAzurePostgreSQLTypeProperties(config LinkedServiceAzurePostgreSQLModel) linkedservices.AzurePostgreSqlLinkedServiceTypeProperties {
	properties := linkedservices.AzurePostgreSqlLinkedServiceTypeProperties{
		Database: pointer.To(any(config.DatabaseName)),
		Server:   pointer.To(any(config.Server)),
		Port:     pointer.To(any(config.Port)),
		Username: pointer.To(any(config.Username)),
	}

	if len(config.KeyVaultPassword) > 0 {
		password := config.KeyVaultPassword[0]
		properties.Password = &linkedservices.AzureKeyVaultSecretReference{
			SecretName: password.SecretName,
			Store: linkedservices.LinkedServiceReference{
				Type:          linkedservices.TypeLinkedServiceReference,
				ReferenceName: password.LinkedServiceName,
			},
		}
	}

	if config.SSLMode != "" {
		properties.SslMode = pointer.To(any(config.SSLMode))
	}
	return properties
}
