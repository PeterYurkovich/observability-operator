package uiplugin

import (
	"fmt"
	"strings"

	"golang.org/x/mod/semver"

	uiv1alpha1 "github.com/rhobs/observability-operator/pkg/apis/uiplugin/v1alpha1"
)

type CompatibilityEntry struct {
	PluginType uiv1alpha1.UIPluginType
	// Minimal OpenShift version supporting this plugin (inclusive).
	MinClusterVersion string
	// Maximal OpenShift version supporting this plugin (exclusive).
	MaxClusterVersion string
	ImageKey          string
	Features          []string
}

var compatibilityMatrix = []CompatibilityEntry{
	{
		PluginType:        uiv1alpha1.TypeDashboards,
		MinClusterVersion: "v4.11",
		MaxClusterVersion: "",
		ImageKey:          "ui-dashboards",
		Features:          []string{},
	},
	{
		PluginType: uiv1alpha1.TypeTroubleshootingPanel,
		// This plugin requires the monitoring-plugin from OpenShift 4.16 (at
		// least) to render the "Troubleshooting Panel" button on the alert
		// details page.
		MinClusterVersion: "v4.16",
		MaxClusterVersion: "",
		ImageKey:          "ui-troubleshooting-panel",
		Features:          []string{},
	},
	{
		PluginType:        uiv1alpha1.TypeDistributedTracing,
		MinClusterVersion: "v4.11",
		MaxClusterVersion: "",
		ImageKey:          "ui-distributed-tracing",
		Features:          []string{},
	},
	{
		PluginType:        uiv1alpha1.TypeLogging,
		MinClusterVersion: "v4.11",
		MaxClusterVersion: "v4.12",
		ImageKey:          "ui-logging",
		Features:          []string{},
	},
	{
		PluginType:        uiv1alpha1.TypeLogging,
		MinClusterVersion: "v4.12",
		MaxClusterVersion: "v4.13",
		ImageKey:          "ui-logging",
		Features: []string{
			"dev-console",
		},
	},
	{
		PluginType:        uiv1alpha1.TypeLogging,
		MinClusterVersion: "v4.13",
		MaxClusterVersion: "v4.14",
		ImageKey:          "ui-logging",
		Features: []string{
			"dev-console",
			"alerts",
		},
	},
	{
		PluginType:        uiv1alpha1.TypeLogging,
		MinClusterVersion: "v4.14",
		MaxClusterVersion: "",
		ImageKey:          "ui-logging",
		Features: []string{
			"dev-console",
			"alerts",
			"dev-alerts",
		},
	},
}

func lookupImageAndFeatures(pluginType uiv1alpha1.UIPluginType, clusterVersion string) (CompatibilityEntry, error) {
	if !strings.HasPrefix(clusterVersion, "v") {
		clusterVersion = "v" + clusterVersion
	}

	// No console plugins are supported before 4.11
	if semver.Compare(clusterVersion, "v4.11") < 0 {
		return CompatibilityEntry{}, fmt.Errorf("dynamic plugins not supported before 4.11")
	}

	for _, entry := range compatibilityMatrix {
		if entry.PluginType != pluginType {
			continue
		}

		canonicalMinClusterVersion := fmt.Sprintf("%s-0", semver.Canonical(entry.MinClusterVersion))
		canonicalMaxClusterVersion := semver.Canonical(entry.MaxClusterVersion)

		if entry.MaxClusterVersion == "" && semver.Compare(clusterVersion, canonicalMinClusterVersion) >= 0 {
			return entry, nil
		}

		if semver.Compare(clusterVersion, canonicalMinClusterVersion) >= 0 && semver.Compare(clusterVersion, canonicalMaxClusterVersion) < 0 {
			return entry, nil
		}
	}

	return CompatibilityEntry{}, fmt.Errorf("plugin %q: no compatible image found for cluster version %q", pluginType, clusterVersion)
}
