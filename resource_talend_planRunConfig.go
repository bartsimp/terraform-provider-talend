package main

import (
	"github.com/bartsimp/talend-rest-go/client/plans_executables"
	"github.com/bartsimp/talend-rest-go/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceTalendPlanRunConfig() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"plan_id": {
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
		Create: resourceTalendPlanRunConfigCreate,
		Read:   resourceTalendPlanRunConfigRead,
		Update: resourceTalendPlanRunConfigUpdate,
		Delete: resourceTalendPlanRunConfigDelete,
	}
}

func resourceTalendPlanRunConfigCreate(d *schema.ResourceData, meta interface{}) error {
	talendClient := meta.(TalendClient)

	planID, body := parsePlanRunConfig(d)

	_, err := talendClient.client.PlansExecutables.ConfigurePlanExecution(
		plans_executables.NewConfigurePlanExecutionParams().WithPlanID(planID).WithBody(&body),
		talendClient.authInfo,
	)
	if err != nil {
		return err
	}

	d.SetId(planID)
	return nil
}

func resourceTalendPlanRunConfigRead(d *schema.ResourceData, meta interface{}) error {
	talendClient := meta.(TalendClient)

	_, err := talendClient.client.PlansExecutables.GetExecutableDetails(
		plans_executables.NewGetExecutableDetailsParams().WithPlanID(d.Id()),
		talendClient.authInfo,
	)
	return err
}

func resourceTalendPlanRunConfigUpdate(d *schema.ResourceData, meta interface{}) error {
	if d.HasChange("plan_id") ||
		d.HasChange("trigger") ||
		d.HasChange("runtime") {

		return resourceTalendPlanRunConfigCreate(d, meta)
	}
	return nil
}

func resourceTalendPlanRunConfigDelete(d *schema.ResourceData, meta interface{}) error {
	talendClient := meta.(TalendClient)

	_, err := talendClient.client.PlansExecutables.DeletePlan(
		plans_executables.NewDeletePlanParams().WithPlanID(d.Id()),
		talendClient.authInfo,
	)
	return err
}

func parsePlanRunConfig(d *schema.ResourceData) (string, models.PlanRunConfig) {
	planID := d.Get("plan_id").(string)

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

	return planID, models.PlanRunConfig{

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
