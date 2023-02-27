package graphendbclient

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

func (client *RestApiClient) CreateEnvironment(ctx context.Context, environmentInfo EnvironmentInfo) (*EnvironmentCreateResult, error) {
	jsonBytes, err := json.MarshalIndent(environmentInfo, "", "  ")
	log.Printf("ola %+v\n", string(jsonBytes))

	response, err := client.ApiClient.R().
		SetBody(environmentInfo).
		SetResult(&EnvironmentCreateResult{}).
		Post("/deployments/environments")
	if err != nil {
		return nil, err
	}
	err = checkResponseAndReturnError(response)
	if err != nil {
		return nil, err
	}
	return response.Result().(*EnvironmentCreateResult), nil
}

// Create vpc peering connection according VpcPeeringInfo provided
func (client *RestApiClient) CreateVpcPeering(ctx context.Context, vpcPeeringInfo VpcPeeringInfo) (*VpcPeeringCreateResult, error) {
	response, err := client.ApiClient.R().
		SetBody(vpcPeeringInfo).
		SetResult(&VpcPeeringCreateResult{}).
		SetPathParams(map[string]string{
			"environmentId": client.EnvironementId,
		}).
		Post("/deployments/environments/{environmentId}/peers")
	if err != nil {
		return nil, err
	}
	err = checkResponseAndReturnError(response)
	if err != nil {
		return nil, err
	}
	return response.Result().(*VpcPeeringCreateResult), nil
}

func (client *RestApiClient) DeleteVpcPeering(ctx context.Context, vpcPeerId string) error {
	response, err := client.ApiClient.R().
		SetPathParams(map[string]string{
			"environmentId": client.EnvironementId,
			"vpcPeerId":     vpcPeerId,
		}).
		Delete("/deployments/environments/{environmentId}/peers/{vpcPeerId}")
	if err != nil {
		return err
	}
	return checkResponseAndReturnError(response)
}

func (client *RestApiClient) DeleteDatabase(ctx context.Context, databaseId string, vendor string) error {
	response, err := client.ApiClient.R().
		SetPathParams(map[string]string{
			"databaseId": databaseId,
			"vendor":     vendor,
		}).
		Delete("/deployments/databases/{vendor}/{databaseId}")
	if err != nil {
		return err
	}
	return checkResponseAndReturnError(response)
}

func (client *RestApiClient) CreateDatabase(ctx context.Context, databaseInfo DatabaseInfo, vendor string) (string, error) {
	databaseInfo.EnvironmentID = client.EnvironementId
	jsonBytes, err := json.MarshalIndent(databaseInfo, "", "  ")
	log.Printf("ola %+v\n", string(jsonBytes))
	response, err := client.ApiClient.R().
		SetBody(databaseInfo).
		SetPathParams(map[string]string{
			"vendor": vendor,
		}).
		SetResult(&DatabaseCreateResult{}).
		Post("/deployments/databases/{vendor}")
	if err != nil {
		return "", err
	}
	err = checkResponseAndReturnError(response)
	if err != nil {
		return "", err
	}
	result := response.Result().(*DatabaseCreateResult)

	if result.Database.Status.State == "paused" {
		_, err = client.FetchDatabaseAsyncStatus(ctx, result.Database.ID, vendor)
	}
	if err != nil {
		return "", err
	}
	return result.Database.ID, nil
}

func (client *RestApiClient) UpdateDatabase(ctx context.Context, databaseId string, databaseInfo DatabaseUpgradeInfo, vendor string) (string, error) {
	response, err := client.ApiClient.R().
		SetBody(databaseInfo).
		SetPathParams(map[string]string{
			"databaseId": databaseId,
			"vendor":     vendor,
		}).
		SetResult(&DatabaseUpdateResult{}).
		Put("/deployments/databases/{vendor}/{databaseId}/plan/change")
	if err != nil {
		return "", err
	}
	err = checkResponseAndReturnError(response)
	if err != nil {
		return "", err
	}
	result := response.Result().(*DatabaseUpdateResult)
	_, err = client.FetchDatabaseAsyncStatus(ctx, databaseId, vendor)
	if err != nil {
		return "", err
	}
	return result.OperationID, nil
}

func (client *RestApiClient) GetUpstreamDatabaseInfo(ctx context.Context, databaseId string, vendor string) (*UpstreamDatabaseInfo, error) {
	response, err := client.ApiClient.R().
		SetPathParams(map[string]string{
			"databaseId": databaseId,
			"vendor":     vendor,
		}).
		SetResult(&UpstreamDatabaseInfo{}).
		Get("/deployments/databases/{vendor}/{databaseId}")
	if response.StatusCode() == 404 {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	err = checkResponseAndReturnError(response)
	if err != nil {
		return nil, err
	}
	return response.Result().(*UpstreamDatabaseInfo), nil
}

func (client *RestApiClient) GetUpstreamDatabasePluginsInfo(ctx context.Context, databaseId string, vendor string) (*PluginListResponse, error) {
	response, err := client.ApiClient.R().
		SetPathParams(map[string]string{
			"databaseId": databaseId,
			"vendor":     vendor,
		}).
		SetResult(&PluginListResponse{}).
		Get("/deployments/databases/{vendor}/{databaseId}/plugins")
	if response.StatusCode() == 404 {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	err = checkResponseAndReturnError(response)
	if err != nil {
		return nil, err
	}
	return response.Result().(*PluginListResponse), nil
}

func (client *RestApiClient) RestartDatabase(ctx context.Context, databaseId string, vendor string) error {
	response, err := client.ApiClient.R().
		SetBody(map[string]interface{}{"reset": true}).
		SetPathParams(map[string]string{
			"databaseId": databaseId,
			"vendor":     vendor,
		}).
		SetResult(&DatabaseRestartResult{}).
		Post("/deployments/databases/{vendor}/{databaseId}/restart")
	if err != nil {
		return err
	}
	err = checkResponseAndReturnError(response)
	if err != nil {
		return err
	}
	result := response.Result().(*DatabaseRestartResult)
	tflog.Debug(ctx, "RESTART API CALL RESULT", map[string]interface{}{
		"StationIds": result.StationIds,
	})
	_, err = client.FetchDatabaseAsyncStatus(ctx, databaseId, vendor)
	if err != nil {
		return err
	}
	return nil
}

func (client *RestApiClient) CreatePlugin(ctx context.Context, databaseId string, pluginInfo PluginInfo, vendor string) (*PluginCreateResult, error) {
	pluginBytes, _ := ioutil.ReadFile(pluginInfo.Url)
	response, err := client.ApiClient.R().
		SetFileReader("file", pluginInfo.Name, bytes.NewReader(pluginBytes)).
		SetFormData(map[string]string{
			"name": pluginInfo.Name,
		}).
		SetPathParams(map[string]string{
			"databaseId": databaseId,
			"vendor":     vendor,
		}).
		SetResult(&PluginCreateResult{}).
		Post("/deployments/databases/{vendor}/{databaseId}/plugins")
	if err != nil {
		return nil, err
	}
	tflog.Debug(ctx, "CREATE PLUGIN IS CALLED", map[string]interface{}{
		"DatabaseId": databaseId,
		"Name":       pluginInfo.Name,
		"Status":     response.Status(),
		"Response":   fmt.Sprintf("%+v", response),
	})
	err = checkResponseAndReturnError(response)
	if err != nil {
		return nil, err
	}
	return response.Result().(*PluginCreateResult), nil
}

func (client *RestApiClient) DeletePlugin(ctx context.Context, databaseId string, pluginId string) error {
	response, err := client.ApiClient.R().
		SetPathParams(map[string]string{
			"databaseId": databaseId,
			"pluginId":   pluginId,
		}).
		Delete("/deployments/databases/{vendor}/{databaseId}/plugins/{pluginId}")
	if err != nil {
		return err
	}
	return checkResponseAndReturnError(response)
}

func (client *RestApiClient) FetchDatabaseAsyncStatus(ctx context.Context, databaseId string, vendor string) (*AsyncDatabaseFetchResult, error) {
	var result *AsyncDatabaseFetchResult
	for {
		response, err := client.ApiClient.R().
			SetPathParams(map[string]string{
				"databaseId": databaseId,
				"vendor":     vendor,
			}).
			SetResult(&AsyncDatabaseFetchResult{}).
			Get("/deployments/databases/{vendor}/{databaseId}/status")
		if err != nil {
			return nil, err
		}
		result = response.Result().(*AsyncDatabaseFetchResult)
		tflog.Debug(ctx, "FETCH ASYNC STATUS API CALL", map[string]interface{}{
			"State":         result.State,
			"NeedsRestart":  result.NeedsRestart,
			"IsPending":     result.IsPending,
			"IsLocked":      result.IsLocked,
			"UnderIncident": result.UnderIncident,
		})
		if !result.IsPending {
			break
		}
		time.Sleep(10 * time.Second)
	}

	if result.State != "running" {
		return nil, fmt.Errorf("Database is not ready. Status result is %+v", result)
	}

	return result, nil
}

func checkResponseAndReturnError(response *resty.Response) error {
	if response.StatusCode() < 200 || response.StatusCode() > 299 {
		return fmt.Errorf("Error happened. Status code is %s and response body is %+v", response.StatusCode(), response)
	}
	return nil
}
