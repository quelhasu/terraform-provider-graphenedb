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
		return "", fmt.Errorf("failed creating the database for operation %s", operationId)
	}
	
	return operationDetail.DatabaseId, nil;
}