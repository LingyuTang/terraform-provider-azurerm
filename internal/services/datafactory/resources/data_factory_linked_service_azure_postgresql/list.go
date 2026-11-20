// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package data_factory_linked_service_azure_postgresql

import (
	"context"
	"fmt"

	"github.com/hashicorp/go-azure-helpers/framework/typehelpers"
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/datafactory/2018-06-01/linkedservices"
	"github.com/hashicorp/terraform-plugin-framework/list"
	"github.com/hashicorp/terraform-plugin-framework/list/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type ListResource struct{}

var _ sdk.FrameworkListWrappedResource = new(ListResource)

func (ListResource) ResourceFunc() *pluginsdk.Resource {
	return RegisterResource()
}

func (ListResource) Metadata(_ context.Context, _ resource.MetadataRequest, response *resource.MetadataResponse) {
	response.TypeName = azurermDataFactoryLinkedServiceAzurePostgreSQLResourceName
}

type LinkedServiceAzurePostgreSQLListModel struct {
	DataFactoryID types.String `tfsdk:"data_factory_id"`
}

func (ListResource) ListResourceConfigSchema(_ context.Context, _ list.ListResourceSchemaRequest, response *list.ListResourceSchemaResponse) {
	response.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"data_factory_id": schema.StringAttribute{
				Required: true,
				Validators: []validator.String{
					typehelpers.WrappedStringValidator{
						Func: linkedservices.ValidateFactoryID,
					},
				},
			},
		},
	}
}

func (ListResource) List(ctx context.Context, request list.ListRequest, stream *list.ListResultsStream, metadata sdk.ResourceMetadata) {
	client := metadata.Client.DataFactory.LinkedServicesClient

	var data LinkedServiceAzurePostgreSQLListModel
	diags := request.Config.Get(ctx, &data)
	if diags.HasError() {
		stream.Results = list.ListResultsStreamDiagnostics(diags)
		return
	}

	dataFactoryId, err := linkedservices.ParseFactoryID(data.DataFactoryID.ValueString())
	if err != nil {
		sdk.SetResponseErrorDiagnostic(stream, "listing azurerm_data_factory_linked_service_azure_postgresql", err)
		return
	}

	resp, err := client.ListByFactoryComplete(ctx, *dataFactoryId)
	if err != nil {
		sdk.SetResponseErrorDiagnostic(stream, fmt.Sprintf("listing `%s`", azurermDataFactoryLinkedServiceAzurePostgreSQLResourceName), err)
		return
	}

	stream.Results = func(push func(list.ListResult) bool) {
		for _, postgreSqlResult := range resp.Items {
			if _, ok := postgreSqlResult.Properties.(linkedservices.AzurePostgreSqlLinkedService); !ok {
				continue
			}

			result := request.NewListResult(ctx)
			result.DisplayName = pointer.From(postgreSqlResult.Name)

			id, err := linkedservices.ParseLinkedServiceID(pointer.From(postgreSqlResult.Id))
			if err != nil {
				sdk.SetErrorDiagnosticAndPushListResult(result, push, "parsing linked service ID", err)
				return
			}

			rd := RegisterResource().Data(&terraform.InstanceState{})
			rd.SetId(id.ID())

			if err := resourceDataFactoryLinkedServiceAzurePostgreSQLFlatten(rd, id, &postgreSqlResult); err != nil {
				sdk.SetErrorDiagnosticAndPushListResult(result, push, fmt.Sprintf("encoding `%s` resource data", azurermDataFactoryLinkedServiceAzurePostgreSQLResourceName), err)
				return
			}

			sdk.EncodeListResult(ctx, rd, &result)
			if result.Diagnostics.HasError() {
				push(result)
				return
			}

			if !push(result) {
				return
			}
		}
	}
}
