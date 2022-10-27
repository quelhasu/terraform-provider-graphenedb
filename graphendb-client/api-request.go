package graphendbclient

import (
	"context"
	"fmt"
	"time"
)

func (client *RestApiClient) CreateVPC(ctx context.Context, vpcInfo VpcInfo) (*VpcCreateResult, error) {
	response, err := client.ApiClient.R().
		SetBody(vpcInfo).
		SetResult(&VpcCreateResult{}).
		Post("/v1/networks")
	if(err != nil){
		return nil, err
	}
	return response.Result().(*VpcCreateResult), nil
}

func (client *RestApiClient) DeleteVPC(ctx context.Context, vpcId string) error {
	_, err := client.ApiClient.R().
		SetPathParams(map[string]string{
			"privateNetworkId": vpcId,
		}).
		Delete("/v1/networks/{privateNetworkId}")
	return err
}


func (client *RestApiClient) CreateDatabase(ctx context.Context, databaseInfo DatabaseInfo) (string, error) {
	response, err := client.ApiClient.R().
		SetBody(databaseInfo).
		SetResult(&DatabaseAsyncOperationResult{}).
		Post("/v1/databases")
	if(err != nil){
		return "", err
	}
	result := response.Result().(*DatabaseAsyncOperationResult)
	asyncOperationInfo, err := client.FetchDatabaseAsyncOperationStatus(result.OperationID)
	if(err != nil) {
		return "", err
	}
	return asyncOperationInfo.DatabaseId, nil
}

func (client *RestApiClient) UpdateDatabase(ctx context.Context, databaseId string, databaseInfo DatabaseUpgradeInfo) (string, error) {
	response, err := client.ApiClient.R().
		SetBody(databaseInfo).
		SetPathParams(map[string]string{
			"databaseId": databaseId,
		}).
		SetResult(&DatabaseAsyncOperationResult{}).
		Put("v1/databases/{databaseId}/upgrade")
	if(err != nil){
		return "", err
	}
	result := response.Result().(*DatabaseAsyncOperationResult)
	asyncOperationInfo, err := client.FetchDatabaseAsyncOperationStatus(result.OperationID)
	if(err != nil) {
		return "", err
	}
	return asyncOperationInfo.DatabaseId, nil
}

func (client *RestApiClient) GetUpstreamDatabaseInfo(ctx context.Context, databaseId string) (*UpstreamDatabaseInfo, error) {
	response, err := client.ApiClient.R().
		SetPathParams(map[string]string{
			"databaseId": databaseId,
		}).
		SetResult(&UpstreamDatabaseInfo{}).
		Get("/v1/databases/{databaseId}")
	if(response != nil && response.StatusCode() == 404) {
		return nil, nil
	}
	if(err != nil) {
		return nil, err
	}
	return response.Result().(*UpstreamDatabaseInfo), nil
}

func  (client *RestApiClient) RestartDatabase(ctx context.Context, databaseId string) error {
	response, err := client.ApiClient.R().
		SetPathParams(map[string]string{
			"databaseId": databaseId,
		}).
		SetResult(&DatabaseRestartResult{}).
		Put("/v1/databases/{databaseId}/restart")
	if(err != nil) {
		return err
	}
	result := response.Result().(*DatabaseRestartResult)
	_, err = client.FetchDatabaseAsyncOperationStatus(result.OperationID)
	if(err != nil) {
		return err
	}
	return nil
}


func (client *RestApiClient) CreatePlugin(ctx context.Context, databaseId string, pluginInfo PluginInfo) (*PluginCreateResult, error) {
	response, err := client.ApiClient.R().
		SetBody(pluginInfo).
		SetPathParams(map[string]string{
			"databaseId": databaseId,
		}).
		SetResult(&PluginCreateResult{}).
		Post("/v1/databases/{databaseId}/plugins")
	if(err != nil){
		return nil, err
	}
	return response.Result().(*PluginCreateResult), nil
}

func (client *RestApiClient) ChangePluginStatus(ctx context.Context, databaseId string, pluginId string, status PluginStatus) error {
	_, err := client.ApiClient.R().
		SetBody(PluginStatusInfo{Status: string(status)}).
		SetPathParams(map[string]string{
			"databaseId": databaseId,
			"pluginId": pluginId,
		}).
		Put("/v1/databases/{databaseId}/plugins/{pluginId}")
	return err
}

func (client *RestApiClient) DeletePlugin(ctx context.Context, databaseId string, pluginId string) error {
	_, err := client.ApiClient.R().
		SetPathParams(map[string]string{
			"databaseId": databaseId,
			"pluginId": pluginId,
		}).
		Delete("/v1/databases/{databaseId}/plugins/{pluginId}")
	return err
}


func (client *RestApiClient) FetchDatabaseAsyncOperationStatus(operationId string) (*AsyncOperationFetchResult, error) {
	var result *AsyncOperationFetchResult;
	for {
		response, err := client.ApiClient.R().
			SetPathParams(map[string]string{
				"operationId": operationId,
			}).
			SetResult(&AsyncOperationFetchResult{}).
			Get("/v1/operations/{operationId}")
		if(err != nil){
			return nil, err
		}
		result = response.Result().(*AsyncOperationFetchResult)
		if(result.Stopped){
			break;
		}
		time.Sleep(10 * time.Second)
	}

	if result.CurrentState != "finished"{
		return nil, fmt.Errorf("async operation is failed. operation id is %s and result is %+v", operationId, result)
	}
	
	return result, nil;
}