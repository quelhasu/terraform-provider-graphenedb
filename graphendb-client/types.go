package graphendbclient

type VpcInfo struct {
	Label     string `json:"label"`
	AwsRegion string `json:"awsRegion"`
	CidrBlock string `json:"cidrBlock"`
}

type VpcCreateResult struct {
	Id string `json:"id"`
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
