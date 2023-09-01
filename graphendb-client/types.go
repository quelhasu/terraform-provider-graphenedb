package graphendbclient

type EnvironmentInfo struct {
	Label          string `json:"label"`
	OrganizationId string `json:"organizationId"`
	GrantType      string `json:"grantType"`
	Region         string `json:"region"`
	Cidr           string `json:"cidr"`
}

type EnvironmentCreateResult struct {
	Id             string        `json:"id"`
	Label          string        `json:"label"`
	OrganizationId string        `json:"organizationId"`
	GrantType      string        `json:"grantType"`
	Region         string        `json:"region"`
	Cidr           string        `json:"cidr"`
	CreatedAt      string        `json:"createdAt"`
	VpcPeers       []VpcPeer     `json:"vpcPeers"`
	NetworkRules   []NetworkRule `json:"networkRules"`
	StoppedAt      string        `json:"stoppedAt"`
}

type VpcPeer struct {
	Id                  string   `json:"id"`
	Label               string   `json:"label"`
	PeeringConnectionId string   `json:"peeringConnectionId"`
	AwsAccountId        string   `json:"awsAccountId"`
	VpcId               string   `json:"vpcId"`
	Cidrs               []string `json:"cidrs"`
	Status              string   `json:"status"`
}

type NetworkRule struct {
	Label     string `json:"label"`
	Id        string `json:"id"`
	CreatedAt string `json:"createdAt"`
	IpRange   string `json:"ipRange"`
}

type VpcPeeringInfo struct {
	Label         string `json:"label"`
	AwsAccountId  string `json:"awsAccountId"`
	VpcId         string `json:"vpcId"`
	PeerVpcRegion string `json:"peerVpcRegion"`
}

type VpcPeeringCreateResult struct {
	ID                  string   `json:"id"`
	Label               string   `json:"label"`
	PeeringConnectionID string   `json:"peeringConnectionId"`
	AWSAccountID        string   `json:"awsAccountId"`
	VPCID               string   `json:"vpcId"`
	CIDRs               []string `json:"cidrs"`
	Status              string   `json:"status"`
}

// type DatabaseInfo struct {
// 	Plan      string `json:"name"`
// 	Version   string `json:"version"`
// 	AwsRegion string `json:"awsRegion"`
// 	Plan      string `json:"plan"`
// 	Vpc       string `json:"privateNetworkId"`
// }

type DatabaseInfo struct {
	EnvironmentID      string        `json:"environmentId"`
	Name               string        `json:"name"`
	Plan               string        `json:"plan"`
	Edition            string        `json:"edition"`
	EnabledExtrasKinds []interface{} `json:"enabledExtrasKinds"`
	Version            string        `json:"version"`
}

type DatabaseUpgradeInfo struct {
	Plan string `json:"targetPlanName"`
}

type DatabaseCreateResult struct {
	Database struct {
		ID                string                   `json:"id"`
		OrganizationID    string                   `json:"organizationId"`
		EnvironmentID     string                   `json:"environmentId"`
		Name              string                   `json:"name"`
		DomainName        string                   `json:"domainName"`
		PrivateDomainName string                   `json:"privateDomainName"`
		HTTPPort          int                      `json:"httpPort"`
		BoltPort          int                      `json:"boltPort"`
		CreatedAt         string                   `json:"createdAt"`
		Status            AsyncDatabaseFetchResult `json:"status"`
		Plan              string                   `json:"plan"`
		Version           struct {
			Number  string `json:"number"`
			Edition string `json:"edition"`
			Vendor  string `json:"vendor"`
		} `json:"version"`
		Nodes []struct {
			StationID string `json:"stationId"`
			NodeType  string `json:"nodeType"`
		} `json:"nodes"`
	} `json:"database"`
	OperationID string `json:"operationId"`
}

type DatabaseUpdateResult struct {
	OperationID string `json:"operationId"`
}

type AsyncOperationFetchResult struct {
	Id                string `json:"id"`
	NextOperationId   string `json:"nextOperationId"`
	Status            string `json:"status"`
	StartedAt         string `json:"startedAt"`
	DurationInSeconds int    `json:"durationInSeconds"`
}

type AsyncDatabaseFetchResult struct {
	State                 string `json:"state"`
	NeedsRestart          bool   `json:"needsRestart"`
	IsPending             bool   `json:"isPending"`
	IsLocked              bool   `json:"isLocked"`
	UnderIncident         bool   `json:"underIncident"`
	LastHealthCheckStatus struct {
		LastCheckSecondsAgo int    `json:"lastCheckSecondsAgo"`
		Status              string `json:"status"`
	} `json:"lastHealthCheckStatus"`
}

type UpstreamDatabaseInfo struct {
	ID                string `json:"id"`
	OrganizationID    string `json:"organizationId"`
	EnvironmentID     string `json:"environmentId"`
	Name              string `json:"name"`
	DomainName        string `json:"domainName"`
	PrivateDomainName string `json:"privateDomainName"`
	HTTPPort          int    `json:"httpPort"`
	BoltPort          int    `json:"boltPort"`
	CreatedAt         string `json:"createdAt"`
	Status            struct {
		State                 string `json:"state"`
		NeedsRestart          bool   `json:"needsRestart"`
		IsPending             bool   `json:"isPending"`
		IsLocked              bool   `json:"isLocked"`
		UnderIncident         bool   `json:"underIncident"`
		LastHealthCheckStatus struct {
			LastCheckSecondsAgo int    `json:"lastCheckSecondsAgo"`
			Status              string `json:"status"`
		} `json:"lastHealthCheckStatus"`
	} `json:"status"`
	Plan    string `json:"plan"`
	Version struct {
		Number  string `json:"number"`
		Edition string `json:"edition"`
		Vendor  string `json:"vendor"`
	} `json:"version"`
	Nodes []struct {
		StationID string `json:"stationId"`
		NodeType  string `json:"nodeType"`
	} `json:"nodes"`
}

type DatabaseRestartResult struct {
	StationIds []string `json:"stationIds"`
	Reset      bool     `json:"reset"`
}

type PluginStatus string

const (
	PluginEnabledStatus  PluginStatus = "enabled"
	PluginDisabledStatus PluginStatus = "disabled"
)

type PluginInfo struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	CreatedAt string `json:"createdAt"`
	Url       string `json:"url"`
}

type PluginListResponse struct {
	Plugins []PluginInfo `json:"plugins"`
}

type PluginStatusInfo struct {
	Status string `json:"status"`
}

type PluginCreateResult struct {
	Id        string `json:"id"`
	Name      string `json:"name"`
	CreatedAt string `json:"createdAt"`
}
