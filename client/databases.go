package client

type DatabasesClient struct {
	ResourceClient
}

type DatabaseSpec struct {
	name      string `json:"name"`
	version   string `json:"version"`
	awsRegion string `json:"awsRegion"`
	plan      string `json:"plan"`
}

type DatabaseDetail struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Version   string `json:"version"`
	AwsRegion string `json:"awsRegion"`
	Plan      string `json:"plan"`
	URI       string `json:"uri"`
}

func (c *AuthenticatedClient) Databases() *DatabasesClient {
	return &DatabasesClient{
		ResourceClient: ResourceClient{
			AuthenticatedClient: c,
			ResourceDescription: "database list",
			ContainerPath:       "/databases/",
			ResourceRootPath:    "/databases",
		}}
}

func (c *DatabasesClient) CreateDatabase(name, version, awsRegion, plan string) (*DatabaseDetail, error) {
	spec := DatabaseSpec{
		name:      name,
		version:   version,
		awsRegion: awsRegion,
		plan:      plan,
	}

	var databaseDetail DatabaseDetail
	if err := c.createResource(&spec, &databaseDetail); err != nil {
		return nil, err
	}

	return c.success(&databaseDetail)
}

func (c *DatabasesClient) success(databaseDetail *DatabaseDetail) (*DatabaseDetail, error) {
	c.unqualify(&databaseDetail.Name)
	return databaseDetail, nil
}
