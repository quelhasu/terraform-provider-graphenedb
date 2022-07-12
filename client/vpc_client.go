package client

type VpcClient struct {
	ResourceClient
}

type VpcSpec struct {
	Label    string `json:"label"`
	AwsRegion  string `json:"awsRegion"`
	CidrBlock    string `json:"cidrBlock"`
}

type VpcDetail struct {
	ID        string `json:"id"`
}

func (c *AuthenticatedClient) NewVpcClient(resourceClients ...ResourceClient) *VpcClient {
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

	return &VpcClient{
		ResourceClient: resourceClient,
	}
}

func (c *VpcClient) CreateVpc(label string, cidr string, region string) (*VpcDetail, error) {
	spec := VpcSpec{
		Label:    label,
		AwsRegion:  region,
		CidrBlock: cidr,
	}

	var vpcDetail VpcDetail
	if err := c.CreateResource(&spec, &vpcDetail); err != nil {
		return nil, err
	}

	return c.success(&vpcDetail)
}

func (c *VpcClient) success(vpcDetail *VpcDetail) (*VpcDetail, error) {
	// 	c.unqualify(&vpcDetail.Name)
	return vpcDetail, nil
}
