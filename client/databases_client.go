package client

type DatabasesClient struct {
	ResourceClient
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