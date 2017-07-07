package graphenedb

import (
	"fmt"

	"github.com/hashicorp/terraform/helper/schema"
)

func resourceDatabase() *schema.Resource {
	return &schema.Resource{
		Create: resourceDatabaseCreate,
		Read:   resourceDatabaseRead,
		Update: resourceDatabaseUpdate,
		Delete: resourceDatabaseDelete,

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"version": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: false,
			},

			"awsRegion": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: false,
			},

			"plan": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: false,
			},
		},
	}
}

func resourceDatabaseCreate(d *schema.ResourceData, meta interface{}) error {
	// cli.Debug.Printf("Resource state: %#v", d.State())

	name, version, awsRegion, plan := getDatabaseResourceData(d)
	client := meta.(*GrapheneDBClient).NewDatabasesClient()
	// cli.Debug.Printf("Resource state: %s %s %s %s %#v", name, version, awsRegion, plan, client)

	_, err := client.CreateDatabase(name, version, awsRegion, plan)
	if err != nil {
		return fmt.Errorf("Error creating database %s: %s", name, err)
	}

	// d.SetId(database.ID)
	return nil
}

func resourceDatabaseRead(d *schema.ResourceData, meta interface{}) error {
	// cli.Debug.Printf("Resource state: %#v", d.State())
	return nil
}

func resourceDatabaseUpdate(d *schema.ResourceData, meta interface{}) error {
	// cli.Debug.Printf("Resource state: %#v", d.State())
	return nil
}

func resourceDatabaseDelete(d *schema.ResourceData, meta interface{}) error {
	// cli.Debug.Printf("Resource state: %#v", d.State())
	return nil
}

func getDatabaseResourceData(d *schema.ResourceData) (string, string, string, string) {
	return d.Get("name").(string),
		d.Get("version").(string),
		d.Get("awsRegion").(string),
		d.Get("plan").(string)
}
