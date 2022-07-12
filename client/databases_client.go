package client

import "log"

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

type DatabaseDetail struct {
	OperationID        string `json:"operationId"`
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

func (c *DatabasesClient) CreateDatabase(name, version, region, plan string, vpc string) (*DatabaseDetail, error) {
	spec := DatabaseSpec{
		Name:    name,
		Version: version,
		AwsRegion:  region,
		Plan:    plan,
		Vpc: vpc,
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
