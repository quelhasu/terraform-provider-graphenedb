package client

type DatabasesClient struct {
	ResourceClient
}

type DatabaseSpec struct {
	name    string `json:"name"`
	version string `json:"version"`
	region  string `json:"region"`
	plan    string `json:"plan"`
	cidr    string `json:"cidr"`
}

type DatabaseDetail struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Version   string `json:"version"`
	AwsRegion string `json:"awsRegion"`
	Plan      string `json:"plan"`
	URI       string `json:"uri"`
}

func (c *AuthenticatedClient) NewDatabasesClient(resourceClients ...ResourceClient) *DatabasesClient {
	var resourceClient ResourceClient

	if len(resourceClients) > 0 {
		resourceClient = resourceClients[0]
	} else {
		resourceClient = &DefaultResourceClient{
			AuthenticatedClient: c,
			ResourceDescription: "database list",
			ResourceRootPath:    "/databases",
		}
	}

	return &DatabasesClient{
		ResourceClient: resourceClient,
	}
}

func (c *DatabasesClient) CreateDatabase(name, version, region, plan string, cidr string) (*DatabaseDetail, error) {
	// 	spec := DatabaseSpec{
	// 		name:      name,
	// 		version:   version,
	// 		region: region,
	// 		plan:      plan,
	// 		cidr: 	cidr
	// 	}

	var databaseDetail DatabaseDetail
	// 	if err := c.createResource(&spec, &databaseDetail); err != nil {
	// 		return nil, err
	// 	}

	return c.success(&databaseDetail)
}

func (c *DatabasesClient) success(databaseDetail *DatabaseDetail) (*DatabaseDetail, error) {
	// 	c.unqualify(&databaseDetail.Name)
	return databaseDetail, nil
}
