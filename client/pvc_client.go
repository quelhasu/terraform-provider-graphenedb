package client

type PvcClient struct {
	ResourceClient
}

type PvcSpec struct {
	Label    string `json:"label"`
	AwsRegion  string `json:"awsRegion"`
	CidrBlock    string `json:"cidrBlock"`
}

type PvcDetail struct {
	ID        string `json:"id"`
}

func (c *AuthenticatedClient) NewPvcClient(resourceClients ...ResourceClient) *PvcClient {
	var resourceClient ResourceClient

	if len(resourceClients) > 0 {
		resourceClient = resourceClients[0]
	} else {
		resourceClient = &DefaultResourceClient{
			AuthenticatedClient: c,
			ResourceDescription: "networks list",
			ResourceRootPath:    "/v1/networks",
		}
	}

	return &PvcClient{
		ResourceClient: resourceClient,
	}
}

func (c *PvcClient) CreatePvc(label string, cidr string, region string) (*PvcDetail, error) {
	spec := PvcSpec{
		Label:    label,
		AwsRegion:  region,
		CidrBlock: cidr,
	}

	var pvcDetail PvcDetail
	if err := c.CreateResource(&spec, &pvcDetail); err != nil {
		return nil, err
	}

	return c.success(&pvcDetail)
}

func (c *PvcClient) success(pvcDetail *PvcDetail) (*PvcDetail, error) {
	// 	c.unqualify(&pvcDetail.Name)
	return pvcDetail, nil
}
