package client

import (
	"context"
	"fmt"
)

type DatabasesClient struct {
	ResourceClient
	ResourcePluginsPath string
}

type DatabaseSpec struct {
	Name    string `json:"name"`
	Version string `json:"version"`
	AwsRegion  string `json:"awsRegion"`
	Plan    string `json:"plan"`
	Vpc			string `json:"privateNetworkId"`
}

type UpgradeDatabaseSpec struct {
	Plan    string `json:"plan"`
}

type DatabaseDetail struct {
	OperationID        string `json:"operation"`
}

type DatabaseUpgradeResult struct {
	OperationID	string `json:"operationId"`
}

type DatabaseRestartDetail struct {
	OperationID	string `json:"operationId"`
}

type PluginSpec struct {
	Name    string `json:"name"`
	Kind string `json:"kind"`
	Url  string `json:"url"`
}

type PluginDetail struct {
	Detail struct {
		Id string `json:"id"`
		Kind string `json:"kind"`
		Enabled bool `json:"enabled"`
		Name string `json:"name"`
	} `json:"plugin"`
}

type StatusPluginDetail struct {}

type StatusPluginSpec struct {
	Status string `json:"status"`
}

type UpstreamDatabasePlanInfo struct {
	PlanType string `json:"type"`
}

type UpstreamDatabaseInfo struct {
	Id string `json:"id"`
	Name string `json:"name"`
	CreatedAt string `json:"createdAt"`
	Version string `json:"version"`
	VersionEdition string `json:"versionEdition"`
	VersionKey string `json:"versionKey"`
	CurrentSize int64 `json:"currentSize"`
	MaxSize int64 `json:"maxSize"`
	Plan UpstreamDatabasePlanInfo `json:"plan"`
	AwsRegion string `json:"awsRegion"`
	PrivateNetworkId string `json:"privateNetworkId"`
	BoltURL string `json:"boltURL"`
	RestUrl string `json:"restUrl"`
	BrowserUrl string `json:"browserUrl"`
	MetricsURL string `json:"metricsURL"`
}

func (c *AuthenticatedClient) NewDatabasesClient(resourceClients ...ResourceClient) *DatabasesClient {
	var resourceClient ResourceClient

	if len(resourceClients) > 0 {
		resourceClient = resourceClients[0]
	} else {
		resourceClient = &DefaultResourceClient{
			AuthenticatedClient: c,
			ResourceDescription: "database list",
			ResourceRootPath:    "/v1/databases",
		}
	}

	return &DatabasesClient{
		ResourceClient: resourceClient,
		ResourcePluginsPath: "%s/plugins",
	}
}

func (c *DatabasesClient) CreateDatabase(name, version, region, plan string, vpc string) (*DatabaseDetail, error) {
	spec := DatabaseSpec{
		Name:    name,
		Version: version,
		AwsRegion:  region,
		Plan:    plan,
		Vpc: vpc,
	}

	var dbCreationDetail DatabaseDetail
	if err := c.CreateResource(&spec, &dbCreationDetail); err != nil {
		return nil, err
	}

	return &dbCreationDetail, nil
}

func (c *DatabasesClient) GetDatabaseInfo(ctx context.Context, id string) (*UpstreamDatabaseInfo, error) {

	var upstreamDatabaseInfo UpstreamDatabaseInfo
	err := c.GetResourceInfo(ctx, fmt.Sprintf("v1/databases/%s", id), &upstreamDatabaseInfo)
	if err != nil {
		apiError, ok := err.(UnexpectedStatusError)
		if(ok && apiError.StatusCode == 404){
				return nil, nil
		}
		return nil, err
	}
	return &upstreamDatabaseInfo, nil
}

func (c *DatabasesClient) UpgradeDatabaseInfo(ctx context.Context, plan string, id string) (*DatabaseUpgradeResult, error) {
	spec := UpgradeDatabaseSpec{
		Plan:    plan,
	}

	var dbCreationDetail DatabaseUpgradeResult
	if err := c.ModifyResource(ctx, fmt.Sprintf("v1/databases/%s/upgrade", id), &spec, &dbCreationDetail); err != nil {
		return nil, err
	}

	return &dbCreationDetail, nil
}

func (c *DatabasesClient) RestartDatabase(database_id string) (*DatabaseRestartDetail, error) {

	var dbCreationDetail DatabaseRestartDetail
	if err := c.RestartResource(database_id, &dbCreationDetail); err != nil {
		return nil, err
	}

	return &dbCreationDetail, nil
}

func (c *DatabasesClient) AddPlugin(name string, kind string, url string, database_id string) (*PluginDetail, error) {
	spec := PluginSpec{
		Name:    name,
		Kind: kind,
		Url: url,
	}

	// extension := fmt.Sprintf("%s/%s", database_id, "plugins")
	var pluginDetail PluginDetail
	var extension = fmt.Sprintf(c.ResourcePluginsPath, database_id)

	if err := c.CreateResourceWithPathExt(extension, &spec, &pluginDetail); err != nil {
		return nil, err
	}

	return &pluginDetail, nil
}

func (c *DatabasesClient) ChangePluginStatus(database_id string, plugin_id string, status string) (*StatusPluginDetail, error) {
	spec := StatusPluginSpec{
		Status:    status,
	}

	var statusDetail StatusPluginDetail
	var plugin_path = fmt.Sprintf(c.ResourcePluginsPath, database_id)
	var extension = fmt.Sprintf("%s/%s", plugin_path, plugin_id) 

	if err := c.ModifyResourceWithPathExt(extension, &spec, &statusDetail); err != nil {
		return nil, err
	}

	return &statusDetail, nil
}