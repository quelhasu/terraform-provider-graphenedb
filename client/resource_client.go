package client

import "fmt"

type ResourceClient interface {
	CreateResource(requestBody interface{}, responseBody interface{}) error
	FetchResource(requestPathId string, responseBody interface{}) error
}

type DefaultResourceClient struct {
	*AuthenticatedClient
	ResourceDescription string
	ResourceRootPath    string
}

func (c *DefaultResourceClient) CreateResource(requestBody interface{}, responseBody interface{}) error {
	request, err := c.newAuthenticatedPostRequest(c.ResourceRootPath, requestBody)
	if err != nil {
		return err
	}

	response, err := c.requestAndCheckStatus(fmt.Sprintf("create %s", c.ResourceDescription), request)
	if err != nil {
		return err
	}

	return unmarshalResponseBody(response, responseBody)
}

func (c *DefaultResourceClient) FetchResource(requestPathId string, responseBody interface{}) error {
	request, err := c.newAuthenticatedGetRequest(c.ResourceRootPath, requestPathId)
	if err != nil {
		return err
	}

	response, err := c.requestAndCheckStatus(fmt.Sprintf("read %s", c.ResourceDescription), request)
	if err != nil {
		return err
	}

	return unmarshalResponseBody(response, responseBody)
}

