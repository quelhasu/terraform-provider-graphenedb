package graphendbclient

import (
	"context"
	"fmt"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

func (client *RestApiClient) CreateVPC(ctx context.Context, vpcInfo VpcInfo) (*VpcCreateResult, error) {
	response, err := client.ApiClient.R().
		SetBody(vpcInfo).
		SetResult(&VpcCreateResult{}).
		Post("/v1/networks")
	if(err != nil){
		return nil, err
	}
	err = checkResponseAndReturnError(response)
	if(err != nil){
		return nil, err
	}
	return response.Result().(*VpcCreateResult), nil
}

func (client *RestApiClient) DeleteVPC(ctx context.Context, vpcId string) error {
	response, err := client.ApiClient.R().
		SetPathParams(map[string]string{
			"privateNetworkId": vpcId,
		}).
		Delete("/v1/networks/{privateNetworkId}")
	if(err != nil){
		return err
	}
	return checkResponseAndReturnError(response)
}


func (client *RestApiClient) CreateDatabase(ctx context.Context, databaseInfo DatabaseInfo) (string, error) {
	response, err := client.ApiClient.R().
		SetBody(databaseInfo).
		SetResult(&DatabaseCreateResult{}).
		Post("/v1/databases")
	if(err != nil){
		return "", err
	}
	err = checkResponseAndReturnError(response)
	if(err != nil){
		return "", err
	}
	result := response.Result().(*DatabaseCreateResult)
	asyncOperationInfo, err := client.FetchDatabaseAsyncOperationStatus(ctx, result.OperationID)
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
		SetResult(&DatabaseUpdateResult{}).
		Put("v1/databases/{databaseId}/upgrade")
	if(err != nil){
		return "", err
	}
	err = checkResponseAndReturnError(response)
	if(err != nil){
		return "", err
	}
	result := response.Result().(*DatabaseUpdateResult)
	asyncOperationInfo, err := client.FetchDatabaseAsyncOperationStatus(ctx, result.OperationID)
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
	if(response.StatusCode() == 404) {
		return nil, nil
	}
	if(err != nil) {
		return nil, err
	}
	err = checkResponseAndReturnError(response)
	if(err != nil){
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
	err = checkResponseAndReturnError(response)
	if(err != nil){
		return err
	}
	result := response.Result().(*DatabaseRestartResult)
	tflog.Debug(ctx, "RESTART API CALL RESULT",  map[string]interface{}{
		"OperationId": result.OperationID,
	})
	_, err = client.FetchDatabaseAsyncOperationStatus(ctx, result.OperationID)
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
	tflog.Debug(ctx, "CREATE PLUGIN IS CALLED",  map[string]interface{}{
		"DatabaseId": databaseId,
		"Kind": pluginInfo.Kind,
		"Url": pluginInfo.Url,
		"Name":  pluginInfo.Name,
		"Status": response.Status(),
		"Response": fmt.Sprintf("%+v", response),
	})
	err = checkResponseAndReturnError(response)
	if(err != nil){
		return nil, err
	}
	return response.Result().(*PluginCreateResult), nil
}

func (client *RestApiClient) ChangePluginStatus(ctx context.Context, databaseId string, pluginId string, status PluginStatus) error {
	response, err := client.ApiClient.R().
		SetBody(PluginStatusInfo{Status: string(status)}).
		SetPathParams(map[string]string{
			"databaseId": databaseId,
			"pluginId": pluginId,
		}).
		Put("/v1/databases/{databaseId}/plugins/{pluginId}")
	tflog.Debug(ctx, "CHANGE STATUS OF PLUGIN API CALL",  map[string]interface{}{
		"databaseId": databaseId,
		"pluginId": pluginId,
		"status": string(status),
		"responseStatus": response.StatusCode(),
	})
	if(err != nil){
		return err
	}
	return checkResponseAndReturnError(response)
}

func (client *RestApiClient) DeletePlugin(ctx context.Context, databaseId string, pluginId string) error {
	response, err := client.ApiClient.R().
		SetPathParams(map[string]string{
			"databaseId": databaseId,
			"pluginId": pluginId,
		}).
		Delete("/v1/databases/{databaseId}/plugins/{pluginId}")
	if(err != nil){
		return err
	}
	return checkResponseAndReturnError(response)
}

func (client *RestApiClient) FetchDatabaseAsyncOperationStatus(ctx context.Context, operationId string) (*AsyncOperationFetchResult, error) {
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
		tflog.Debug(ctx, "FETCH ASYNC OPERATION STATUS API CALL",  map[string]interface{}{
			"Id": result.Id,
			"DatabaseId": result.DatabaseId,
			"Description": result.Description,
			"CurrentState": result.CurrentState,
			"Stopped": result.Stopped,
		})
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

func checkResponseAndReturnError(response *resty.Response) error {
	if response.StatusCode() < 200 || response.StatusCode() > 299 {
		return fmt.Errorf("Error happened. Status code is %s and response body is %+v", response.StatusCode(), response)
	}
	return nil
}