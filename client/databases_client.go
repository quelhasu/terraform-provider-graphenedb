package client

import "log"

type DatabasesClient struct {
	ResourceClient
}

type DatabaseSpec struct {
	name    string `json:"name"`
	version string `json:"version"`
	awsRegion  string `json:"awsRegion"`
	plan    string `json:"plan"`
}

type DatabaseDetail struct {
	OperationID        string `json:"id"`
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

func (c *DatabasesClient) CreateDatabase(name, version, region, plan string) (*DatabaseDetail, error) {
	spec := DatabaseSpec{
		name:    name,
		version: version,
		awsRegion:  region,
		plan:    plan,
	}

	var databaseDetail DatabaseDetail
	if err := c.CreateResource(&spec, &databaseDetail); err != nil {
		log.Printf("DatabaseDetail: ", err)
		return nil, err
	}

	return c.success(&databaseDetail)
}

func (c *DatabasesClient) success(databaseDetail *DatabaseDetail) (*DatabaseDetail, error) {
	// 	c.unqualify(&databaseDetail.Name)
	return databaseDetail, nil
}
