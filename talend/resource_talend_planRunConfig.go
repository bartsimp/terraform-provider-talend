package talend

import (
	talendRest "github.com/bartsimp/talend-rest-go"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceTalendPlanRunConfig() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"json_request": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
		},
		Create: resourceTalendPlanRunConfigCreate,
		Read:   resourceTalendPlanRunConfigRead,
		Update: resourceTalendPlanRunConfigUpdate,
		Delete: resourceTalendPlanRunConfigDelete,
	}
}

func resourceTalendPlanRunConfigCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*talendRest.Client)
	jsonRequest := d.Get("json_request").(string)

	planRunConfig, err := client.UpdatePlanRunConfigFromPlainJson(jsonRequest)
	if err != nil {
		return err
	}

	d.SetId(planRunConfig.PlanId)
	return resourceTalendPlanRunConfigRead(d, meta)
}

func resourceTalendPlanRunConfigRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*talendRest.Client)
	id := d.Id()

	planRunConfig, err := client.GetPlanRunConfigByPlanId(id)
	if err != nil {
		return err
	}

	d.SetId(planRunConfig.PlanId)
	//	d.Set("name", planRunConfig.Trigger.Webhook.Name)

	return nil
}

func resourceTalendPlanRunConfigUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*talendRest.Client)

	jsonRequest := d.Get("json_request").(string)

	if d.HasChange("json_request") {
		planRunConfig, err := client.UpdatePlanRunConfigFromPlainJson(jsonRequest)
		if err != nil {
			return err
		}

		d.Set("id", planRunConfig.PlanId)
		//		d.Set("name", planRunConfig.Trigger.Webhook.Name)

	}

	return nil
}

func resourceTalendPlanRunConfigDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*talendRest.Client)

	err := client.DeletePlanRunConfigByPlanId(d.Id())
	if err != nil {
		return err
	}

	return nil
}
