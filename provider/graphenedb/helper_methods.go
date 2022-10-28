package graphenedb

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func AttributesToResourceData(apiAttributes map[string]interface{}, d *schema.ResourceData) error {
	for attributeName, attributeValue := range apiAttributes {
		if err := d.Set(attributeName, attributeValue); err != nil {
			return fmt.Errorf("error setting %s: %w", attributeName, err)
		}
	}
	return nil
}
