package talend

import (
	talendRest "github.com/bartsimp/talend-rest-go"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceTalendPlan() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"json_request": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
		},
		Create: resourceTalendPlanCreate,
		Read:   resourceTalendPlanRead,
		Update: resourceTalendPlanUpdate,
		Delete: resourceTalendPlanDelete,
	}
}

func resourceTalendPlanCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*talendRest.Client)
	jsonRequest := d.Get("json_request").(string)

	plan, err := client.CreatePlanFromPlainJSON(jsonRequest)
	if err != nil {
		return err
	}

	d.SetId(plan.Id)
	return resourceTalendPlanRead(d, meta)
}

func resourceTalendPlanRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*talendRest.Client)
	id := d.Id()

	plan, err := client.GetPlanByID(id)
	if err != nil {
		return err
	}

	d.Set("id", plan.Executable)
	d.Set("name", plan.Name)

	return nil
}

func resourceTalendPlanUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*talendRest.Client)

	jsonRequest := d.Get("json_request").(string)

	if d.HasChange("json_request") {
		plan, err := client.UpdatePlanFromPlainJSON(d.Id(), jsonRequest)
		if err != nil {
			return err
		}

		d.Set("id", plan.Id)
		d.Set("name", plan.Name)

	}

	return nil
}

func resourceTalendPlanDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*talendRest.Client)

	err := client.DeletePlan(d.Id())
	if err != nil {
		return err
	}

	return nil
}
