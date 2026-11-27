// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package kubernetes_fleet_update_run

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-sdk/resource-manager/containerservice/2025-03-01/updateruns"
	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

var (
	_ sdk.Resource                  = Resource{}
	_ sdk.ResourceWithUpdate        = Resource{}
	_ sdk.ResourceWithCustomizeDiff = Resource{}
)

type Resource struct{}

func (r Resource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return updateruns.ValidateUpdateRunID
}

func (r Resource) ResourceType() string {
	return "azurerm_kubernetes_fleet_update_run"
}

func (r Resource) Description() string {
	return "Manages a Kubernetes Fleet Update Run."
}

func (r Resource) CustomizeDiff() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			rawConfig := metadata.ResourceDiff.GetRawConfig()
			if !rawConfig.IsKnown() || rawConfig.IsNull() {
				return nil
			}

			return validateKubernetesVersion(rawConfig)
		},
	}
}

func validateKubernetesVersion(rawConfig cty.Value) error {
	managedClusterUpdate := rawConfig.GetAttr("managed_cluster_update")
	if !managedClusterUpdate.IsKnown() {
		return nil
	}

	upgrade := managedClusterUpdate.AsValueSlice()[0].GetAttr("upgrade")
	if !upgrade.IsKnown() {
		return nil
	}

	upgradeConfig := upgrade.AsValueSlice()[0]
	upgradeType := upgradeConfig.GetAttr("type")
	if !upgradeType.IsKnown() {
		return nil
	}

	upgradeTypeEnum := updateruns.ManagedClusterUpgradeType(upgradeType.AsString())
	if upgradeTypeEnum != updateruns.ManagedClusterUpgradeTypeFull && upgradeTypeEnum != updateruns.ManagedClusterUpgradeTypeControlPlaneOnly {
		return nil
	}

	kubernetesVersion := upgradeConfig.GetAttr("kubernetes_version")
	if !kubernetesVersion.IsKnown() {
		return nil
	}
	if kubernetesVersion.IsNull() {
		return fmt.Errorf("`managed_cluster_update.0.upgrade.0.kubernetes_version` is required when `managed_cluster_update.0.upgrade.0.type` is `%s`", upgradeType.AsString())
	}

	return nil
}
