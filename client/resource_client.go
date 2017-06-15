package client

import "fmt"

type ResourceClient struct {
	*AuthenticatedClient
	ResourceDescription string
	ContainerPath       string
	ResourceRootPath    string
}

func (c *ResourceClient) createResource(requestBody interface{}, responseBody interface{}) error {
	request, err := c.newAuthenticatedPostRequest(c.ContainerPath, requestBody)
	if err != nil {
		return err
	}

	response, err := c.requestAndCheckStatus(fmt.Sprintf("create %s", c.ResourceDescription), request)
	if err != nil {
		return err
	}

	return unmarshalResponseBody(response, responseBody)
}
