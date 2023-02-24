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
	Name      string `json:"name"`
	Version   string `json:"version"`
	AwsRegion string `json:"awsRegion"`
	Plan      string `json:"plan"`
	Vpc       string `json:"privateNetworkId"`
}

type DatabaseUpgradeInfo struct {
	Plan string `json:"plan"`
}

type DatabaseCreateResult struct {
	OperationID string `json:"operation"`
}

type DatabaseUpdateResult struct {
	OperationID string `json:"operationId"`
}

type AsyncOperationFetchResult struct {
	Id           string `json:"id"`
	DatabaseId   string `json:"databaseId"`
	Description  string `json:"description"`
	CurrentState string `json:"currentState"`
	Stopped      bool   `json:"stopped"`
}

type UpstreamDatabaseInfo struct {
	Id             string `json:"id"`
	Name           string `json:"name"`
	CreatedAt      string `json:"createdAt"`
	Version        string `json:"version"`
	VersionEdition string `json:"versionEdition"`
	VersionKey     string `json:"versionKey"`
	CurrentSize    int64  `json:"currentSize"`
	MaxSize        int64  `json:"maxSize"`
	Plan           struct {
		PlanType string `json:"type"`
	} `json:"plan"`
	AwsRegion        string `json:"awsRegion"`
	PrivateNetworkId string `json:"privateNetworkId"`
	BoltURL          string `json:"boltURL"`
	RestUrl          string `json:"restUrl"`
	BrowserUrl       string `json:"browserUrl"`
	MetricsURL       string `json:"metricsURL"`
	Plugins          []struct {
		Id        string `json:"id"`
		Name      string `json:"name"`
		CreatedAt string `json:"createdAt"`
		Enabled   bool   `json:"enabled"`
		Type      string `json:"type"`
	} `json:"plugins"`
}

type DatabaseRestartResult struct {
	OperationID string `json:"operationId"`
}

type PluginStatus string

const (
	PluginEnabledStatus  PluginStatus = "enabled"
	PluginDisabledStatus PluginStatus = "disabled"
)

type PluginInfo struct {
	Name string `json:"name"`
	Kind string `json:"kind"`
	Url  string `json:"url"`
}

type PluginStatusInfo struct {
	Status string `json:"status"`
}

type PluginCreateResult struct {
	Detail struct {
		Id      string `json:"id"`
		Kind    string `json:"kind"`
		Enabled bool   `json:"enabled"`
		Name    string `json:"name"`
	} `json:"plugin"`
}
