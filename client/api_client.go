package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
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

	log.Printf("Request Header:", request.Header)
	log.Printf("Request :", request)
	return request, nil
}

func (c *apiClient) newRequest(method, path string, body interface{}) (*http.Request, error) {
	urlPath, err := url.Parse(path)
	if err != nil {
		return nil, err
	}
	log.Printf("Request Path: ", path)
	log.Printf("Request Body : ", body)
	req, err := http.NewRequest(
		method,
		c.apiEndpoint.ResolveReference(urlPath).String(),
		marshalToReader(body),
	)

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
