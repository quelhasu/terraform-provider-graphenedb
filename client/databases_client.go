package client

import "fmt"

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
type DatabaseCreationDetail struct {
	OperationID        string `json:"operation"`
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

func (c *DatabasesClient) CreateDatabase(name, version, region, plan string, vpc string) (*DatabaseCreationDetail, error) {
	spec := DatabaseSpec{
		Name:    name,
		Version: version,
		AwsRegion:  region,
		Plan:    plan,
		Vpc: vpc,
	}

	var dbCreationDetail DatabaseCreationDetail
	if err := c.CreateResource(&spec, &dbCreationDetail); err != nil {
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