package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
)

type apiClient struct {
	apiEndpoint *url.URL
	httpClient  *http.Client
}

type UnexpectedStatusError struct {
	Description string
	Status      string
	StatusCode  int
	Body        string
}

func (c *apiClient) requestAndCheckStatus(description string, req *http.Request) (*http.Response, error) {
	rsp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	if 300 > rsp.StatusCode && rsp.StatusCode >= 200 {
		return rsp, nil
	}

	return nil, unexpectedStatusError(description, rsp)
}

func (c *apiClient) newPostRequest(path string, body interface{}, api_key string) (*http.Request, error) {
	request, err := c.newRequest("POST", path, body)
	if err != nil {
		return nil, err
	}

	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("api_key", api_key)

	return request, nil
}

func (c *apiClient) newPutRequest(path string, body interface{}, api_key string) (*http.Request, error) {
	request, err := c.newRequest("PUT", path, body)
	if err != nil {
		return nil, err
	}

	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("api_key", api_key)

	return request, nil
}

func (c *apiClient) newGetRequest(path string, api_key string) (*http.Request, error) {
	request, err := c.newRequest("GET", path, nil)
	if err != nil {
		return nil, err
	}

	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("api_key", api_key)

	return request, nil
}

func (c *apiClient) newRequest(method, path string, body interface{}) (*http.Request, error) {
	urlPath, err := url.Parse(path)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest(
		method,
		c.apiEndpoint.ResolveReference(urlPath).String(),
		marshalToReader(body),
	)

	requestDump, err := httputil.DumpRequest(req, true)
	if err != nil {
		log.Println(err)
	}
	log.Println(string(requestDump))

	if err != nil {
		return nil, err
	}

	return req, nil
}

func marshalToReader(body interface{}) io.Reader {
	if body == nil {
		return nil
	}
	bodyData, err := json.Marshal(body)
	if err != nil {
		log.Panic(err)
	}
	return bytes.NewReader(bodyData)
}

func unmarshalResponseBody(response *http.Response, meta interface{}) error {
	buf := new(bytes.Buffer)
	buf.ReadFrom(response.Body)
	return json.Unmarshal(buf.Bytes(), meta)
}

func unexpectedStatusError(description string, rsp *http.Response) error {
	var bodyString string
	if rsp.Body == nil {
		bodyString = "<empty body>"
	} else {
		buf := new(bytes.Buffer)
		buf.ReadFrom(rsp.Body)
		bodyString = buf.String()
	}

	return UnexpectedStatusError{
		Description: description,
		Status:      rsp.Status,
		StatusCode:  rsp.StatusCode,
		Body:        bodyString,
	}
}

func (e UnexpectedStatusError) Error() string {
	return fmt.Sprintf("Unable to %s: %s\n%s", e.Description, e.Status, e.Body)
}
