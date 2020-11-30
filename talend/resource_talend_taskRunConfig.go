package talend

import (
	talendRest "github.com/bartsimp/talend-rest-go"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceTalendTaskRunConfig() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"json_request": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
		},
		Create: resourceTalendTaskRunConfigCreate,
		Read:   resourceTalendTaskRunConfigRead,
		Update: resourceTalendTaskRunConfigUpdate,
		Delete: resourceTalendTaskRunConfigDelete,
	}
}

func resourceTalendTaskRunConfigCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*talendRest.Client)
	jsonRequest := d.Get("json_request").(string)

	taskRunConfig, err := client.UpdateTaskRunConfigFromPlainJson(jsonRequest)
	if err != nil {
		return err
	}

	d.SetId(taskRunConfig.TaskId)
	return resourceTalendTaskRunConfigRead(d, meta)
}

func resourceTalendTaskRunConfigRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*talendRest.Client)
	id := d.Id()

	taskRunConfig, err := client.GetTaskRunConfigByTaskId(id)
	if err != nil {
		return err
	}

	d.SetId(taskRunConfig.TaskId)
	//	d.Set("name", taskRunConfig.Trigger.Webhook.Name)

	return nil
}

func resourceTalendTaskRunConfigUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*talendRest.Client)

	jsonRequest := d.Get("json_request").(string)

	if d.HasChange("json_request") {
		taskRunConfig, err := client.UpdateTaskRunConfigFromPlainJson(jsonRequest)
		if err != nil {
			return err
		}

		d.Set("id", taskRunConfig.TaskId)
		//		d.Set("name", taskRunConfig.Trigger.Webhook.Name)

	}

	return nil
}

func resourceTalendTaskRunConfigDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*talendRest.Client)

	err := client.DeleteTaskRunConfigByTaskId(d.Id())
	if err != nil {
		return err
	}

	return nil
}
