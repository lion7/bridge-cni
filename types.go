package main

// partially copied from https://github.com/containernetworking/cni/blob/main/pkg/types/types.go

// NetConfList describes an ordered list of networks.
type NetConfList struct {
	CNIVersion string `json:"cniVersion,omitempty"`

	Name         string        `json:"name,omitempty"`
	DisableCheck bool          `json:"disableCheck,omitempty"`
	DisableGC    bool          `json:"disableGC,omitempty"`
	Plugins      []*PluginConf `json:"plugins,omitempty"`
}

// PluginConf describes a plugin configuration for a specific network.
type PluginConf struct {
	CNIVersion string `json:"cniVersion,omitempty"`

	Name             string          `json:"name,omitempty"`
	Type             string          `json:"type,omitempty"`
	Capabilities     map[string]bool `json:"capabilities,omitempty"`
	IPAM             IPAM            `json:"ipam,omitempty"`
	DNS              DNS             `json:"dns,omitempty"`
	IsDefaultGateway bool            `json:"isDefaultGateway,omitempty"` // added bridge specific field
}

type IPAM struct {
	Type   string `json:"type,omitempty"`
	Subnet string `json:"subnet,omitempty"` // added host-local specific field
}

// DNS contains values interesting for DNS resolvers
type DNS struct {
	Nameservers []string `json:"nameservers,omitempty"`
	Domain      string   `json:"domain,omitempty"`
	Search      []string `json:"search,omitempty"`
	Options     []string `json:"options,omitempty"`
}
