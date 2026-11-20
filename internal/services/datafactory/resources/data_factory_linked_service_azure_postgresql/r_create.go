// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package data_factory_linked_service_azure_postgresql

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/datafactory/2018-06-01/linkedservices"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/datafactory/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceDataFactoryLinkedServiceAzurePostgreSQLCreateUpdate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).DataFactory.LinkedServicesClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	dataFactoryId, err := linkedservices.ParseFactoryID(d.Get("data_factory_id").(string))
	if err != nil {
		return err
	}

	id := linkedservices.NewLinkedServiceID(meta.(*clients.Client).Account.SubscriptionId, dataFactoryId.ResourceGroupName, dataFactoryId.FactoryName, d.Get("name").(string))
	if d.IsNewResource() {
		if !meta.(*clients.Client).Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
			existing, err := client.Get(ctx, id, linkedservices.DefaultGetOperationOptions())
			if err != nil && !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
			}
			if !response.WasNotFound(existing.HttpResponse) {
				return tf.ImportAsExistsError(azurermDataFactoryLinkedServiceAzurePostgreSQLResourceName, id.ID())
			}
		}
	} else {
		existingId, err := linkedservices.ParseLinkedServiceID(d.Id())
		if err != nil {
			return err
		}
		id = *existingId
	}

	parameters := linkedservices.LinkedServiceResource{
		Properties: linkedservices.AzurePostgreSqlLinkedService{
			Annotations:    helpers.ExpandTypedLinkedServiceAnnotations(d.Get("annotations").([]any)),
			ConnectVia:     helpers.ExpandTypedLinkedServiceIntegrationRuntimeName(d.Get("integration_runtime_name").(string)),
			Description:    pointer.To(d.Get("description").(string)),
			Parameters:     helpers.ExpandTypedLinkedServiceParameters(d.Get("parameters").(map[string]any)),
			TypeProperties: expandLinkedServiceAzurePostgreSQLTypeProperties(d),
			Version:        pointer.To("2.0"),
		},
	}

	if _, err := client.CreateOrUpdate(ctx, id, parameters, linkedservices.DefaultCreateOrUpdateOperationOptions()); err != nil {
		return fmt.Errorf("creating/updating %s: %+v", id, err)
	}

	d.SetId(id.ID())
	if err := pluginsdk.SetResourceIdentityData(d, &id); err != nil {
		return err
	}
	return resourceDataFactoryLinkedServiceAzurePostgreSQLRead(d, meta)
}

func expandLinkedServiceAzurePostgreSQLTypeProperties(d *pluginsdk.ResourceData) linkedservices.AzurePostgreSqlLinkedServiceTypeProperties {
	properties := linkedservices.AzurePostgreSqlLinkedServiceTypeProperties{
		Database: pointer.To(any(d.Get("database_name").(string))),
		Server:   pointer.To(any(d.Get("server").(string))),
		Port:     pointer.To(any(int64(d.Get("port").(int)))),
		Username: pointer.To(any(d.Get("username").(string))),
		Password: helpers.ExpandTypedLinkedServiceKeyVaultPassword(d.Get("key_vault_password").([]any)),
	}

	if sslMode := d.Get("ssl_mode").(string); sslMode != "" {
		properties.SslMode = pointer.To(any(sslMode))
	}
	return properties
}
