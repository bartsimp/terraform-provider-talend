package main

import (
	"fmt"

	"github.com/bartsimp/talend-rest-go/client/tasks"
	"github.com/bartsimp/talend-rest-go/models"
	"github.com/bartsimp/talend-rest-go/utils"
	"github.com/go-openapi/runtime"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceTalendTaskRunConfig() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"task_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"trigger": {
				Type:     schema.TypeSet,
				Required: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:     schema.TypeString,
							Required: true,
						},
						"start_date": {
							Type:     schema.TypeString,
							Required: true,
						},
						"time_zone": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
			"runtime": {
				Type:     schema.TypeSet,
				Required: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
		},
		Create: resourceTalendTaskRunConfigCreate,
		Read:   resourceTalendTaskRunConfigRead,
		Update: resourceTalendTaskRunConfigUpdate,
		Delete: resourceTalendTaskRunConfigDelete,
	}
}

func resourceTalendTaskRunConfigCreate(d *schema.ResourceData, meta interface{}) error {
	talendClient := meta.(TalendClient)

	taskId, body := parseTaskRunConfig(d)

	_, err := talendClient.client.Tasks.ConfigureTaskExecution(
		tasks.NewConfigureTaskExecutionParams().WithTaskID(taskId).WithBody(&body),
		func(co *runtime.ClientOperation) {
			co.AuthInfo = talendClient.authInfo
		},
	)
	if err != nil {
		switch err := err.(type) {
		case *tasks.ConfigureTaskExecutionBadRequest:
			return fmt.Errorf("%s", utils.UnmarshalErrorResponse(err.GetPayload()))
		case *tasks.ConfigureTaskExecutionUnauthorized:
			return fmt.Errorf("unauthorized %s", utils.UnmarshalErrorResponse(err.GetPayload()))
		}
		return err
	}

	d.SetId(taskId)
	return nil
}

func resourceTalendTaskRunConfigRead(d *schema.ResourceData, meta interface{}) error {
	talendClient := meta.(TalendClient)

	_, err := talendClient.client.Tasks.GetTaskConfiguration(
		tasks.NewGetTaskConfigurationParams().WithTaskID(d.Id()),
		func(co *runtime.ClientOperation) {
			co.AuthInfo = talendClient.authInfo
		},
	)
	return err
}

func resourceTalendTaskRunConfigUpdate(d *schema.ResourceData, meta interface{}) error {
	if d.HasChange("task_id") ||
		d.HasChange("trigger") ||
		d.HasChange("runtime") {

		return resourceTalendPlanRunConfigCreate(d, meta)
	}
	return nil
}

func resourceTalendTaskRunConfigDelete(d *schema.ResourceData, meta interface{}) error {
	talendClient := meta.(TalendClient)

	_, err := talendClient.client.Tasks.StopSchedule(
		tasks.NewStopScheduleParams().WithTaskID(d.Id()),
		talendClient.authInfo,
	)
	return err
}

func parseTaskRunConfig(d *schema.ResourceData) (string, models.TaskRunConfig) {
	taskID := d.Get("task_id").(string)

	setTrigger := d.Get("trigger").(*schema.Set)
	trigger0 := setTrigger.List()[0]
	trigger := trigger0.(map[string]interface{})
	triggerType := trigger["type"].(string)
	startDate := trigger["start_date"].(string)
	timeZone := trigger["time_zone"].(string)

	setRuntime := d.Get("runtime").(*schema.Set)
	runtime0 := setRuntime.List()[0]
	runtime := runtime0.(map[string]interface{})
	runtimeType := runtime["type"].(string)

	return taskID, models.TaskRunConfig{
		Trigger: &models.Trigger{
			Type:      &triggerType,
			StartDate: &startDate,
			TimeZone:  &timeZone,
		},
		Runtime: &models.Runtime{
			Type: runtimeType,
		},
	}
}
