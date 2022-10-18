package graphenedb

import (
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/quelhasu/terraform-provider-graphenedb/client"
)

func AttributesToResourceData(apiAttributes map[string]interface{}, d *schema.ResourceData) error {
	for attributeName, attributeValue := range apiAttributes {
		if err := d.Set(attributeName, attributeValue); err != nil {
			return fmt.Errorf("error setting %s: %w", attributeName, err)
		}
	}
	return nil
}

func FetchDatabaseAsyncOperationStatus(operationId string, client *client.OperationsClient) (string, error) {
	operationDetail, err_op := client.FetchOperationDetail(operationId)

	for !operationDetail.Stopped && err_op == nil  {
		if err_op != nil {
			return "", fmt.Errorf("error fetching operation %s: %s", operationId, err_op)
		}
		time.Sleep(10 * time.Second)
		operationDetail, err_op = client.FetchOperationDetail(operationId)
	}

	if operationDetail.CurrentState != "finished"{
		return "", fmt.Errorf("async operation is failed. operation id is %s and result is %+v", operationId, operationDetail)
	}
	
	return operationDetail.DatabaseId, nil;
}

func RestartDatabase(databaseId string, client *GrapheneDBClient) error {
	
	databaseClient := client.NewDatabasesClient()
	databaseRestartResult, err := databaseClient.RestartDatabase(databaseId)
	if err != nil {
		return err
	}

	_, err = FetchDatabaseAsyncOperationStatus(databaseRestartResult.OperationID, client.NewOperationsClient())
	if err != nil {
		return err
	}
	return nil
}