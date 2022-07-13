package client

type OperationsClient struct {
	ResourceClient
}

type OperationDetail struct {
	Id        string `json:"id"`
	DatabaseId        string `json:"databaseId"`
	Description        string `json:"description"`
	CurrentState        string `json:"currentState"`
	Stopped        bool `json:"stopped"`
}

func (c *AuthenticatedClient) NewOperationsClient(resourceClients ...ResourceClient) *OperationsClient {
	var resourceClient ResourceClient

	if len(resourceClients) > 0 {
		resourceClient = resourceClients[0]
	} else {
		resourceClient = &DefaultResourceClient{
			AuthenticatedClient: c,
			ResourceDescription: "operation list",
			ResourceRootPath:    "/v1/operations",
		}
	}

	return &OperationsClient{
		ResourceClient: resourceClient,
	}
}

func (c *OperationsClient) FetchOperationDetail(operationId string) (*OperationDetail, error) {

	var operationDetail OperationDetail
	if err := c.FetchResource(operationId, &operationDetail); err != nil {
		return nil, err
	}

	return &operationDetail, nil
}
