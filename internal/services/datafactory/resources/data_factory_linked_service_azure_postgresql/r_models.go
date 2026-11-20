// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package data_factory_linked_service_azure_postgresql

type LinkedServiceAzurePostgreSQLModel struct {
	Name                   string             `tfschema:"name"`
	DataFactoryID          string             `tfschema:"data_factory_id"`
	Annotations            []string           `tfschema:"annotations"`
	DatabaseName           string             `tfschema:"database_name"`
	Description            string             `tfschema:"description"`
	IntegrationRuntimeName string             `tfschema:"integration_runtime_name"`
	KeyVaultPassword       []KeyVaultPassword `tfschema:"key_vault_password"`
	Parameters             map[string]any     `tfschema:"parameters"`
	Port                   int64              `tfschema:"port"`
	Server                 string             `tfschema:"server"`
	SSLMode                string             `tfschema:"ssl_mode"`
	Username               string             `tfschema:"username"`
}

type KeyVaultPassword struct {
	LinkedServiceName string `tfschema:"linked_service_name"`
	SecretName        string `tfschema:"secret_name"`
}

func (r Resource) ModelObject() any {
	return &LinkedServiceAzurePostgreSQLModel{}
}
