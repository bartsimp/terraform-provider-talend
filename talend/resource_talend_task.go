package talend

import (
	talendRest "github.com/bartsimp/talend-rest-go"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceTalendTask() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"json_request": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
		},
		Create: resourceTalendTaskCreate,
		Read:   resourceTalendTaskRead,
		Update: resourceTalendTaskUpdate,
		Delete: resourceTalendTaskDelete,
	}
}

func resourceTalendTaskCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*talendRest.Client)
	jsonRequest := d.Get("json_request").(string)

	task, err := client.CreateTaskFromPlainJson(jsonRequest)
	if err != nil {
		return err
	}

	d.SetId(task.Id)
	return resourceTalendTaskRead(d, meta)
}

func resourceTalendTaskRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*talendRest.Client)

	task, err := client.GetTaskById(d.Id())
	if err != nil {
		return err
	}

	d.Set("id", task.Id)
	d.Set("name", task.Name)

	return nil
}

func resourceTalendTaskUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*talendRest.Client)

	jsonRequest := d.Get("json_request").(string)

	if d.HasChange("json_request") {
		task, err := client.UpdateTaskFromPlainJson(d.Id(), jsonRequest)
		if err != nil {
			return err
		}

		d.Set("id", task.Id)
		d.Set("name", task.Name)

	}

	return nil
}

func resourceTalendTaskDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*talendRest.Client)

	err := client.DeleteTask(d.Id())
	if err != nil {
		return err
	}

	return nil
}
