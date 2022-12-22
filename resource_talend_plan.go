package main

import (
	"fmt"

	"github.com/bartsimp/talend-rest-go/client/plans_executables"
	"github.com/bartsimp/talend-rest-go/models"
	"github.com/bartsimp/talend-rest-go/utils"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceTalendPlan() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"workspace_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"steps": {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"condition": {
							Type:     schema.TypeString,
							Required: true,
						},
						"task_ids": {
							Type:     schema.TypeList,
							Required: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
					},
				},
			},
		},
		Create: resourceTalendPlanCreate,
		Read:   resourceTalendPlanRead,
		Update: resourceTalendPlanUpdate,
		Delete: resourceTalendPlanDelete,
	}
}

func resourceTalendPlanCreate(d *schema.ResourceData, meta interface{}) error {
	talendClient := meta.(TalendClient)

	planRequest := parsePlanRequest(d)

	createPlanCreated, err := talendClient.client.PlansExecutables.CreatePlan(
		plans_executables.NewCreatePlanParams().WithBody(&planRequest),
		talendClient.authInfo,
	)
	if err != nil {
		switch err := err.(type) {
		case *plans_executables.GetAvailablePlansBadRequest:
		case *plans_executables.CreatePlanBadRequest:
			return fmt.Errorf("%s", utils.UnmarshalErrorResponse(err.GetPayload()))
		}
		return err
	}

	d.SetId(*createPlanCreated.GetPayload().ID)
	return resourceTalendPlanRead(d, meta)
}

func resourceTalendPlanRead(d *schema.ResourceData, meta interface{}) error {
	talendClient := meta.(TalendClient)

	_, err := talendClient.client.PlansExecutables.GetExecutableDetails(
		plans_executables.NewGetExecutableDetailsParams().WithPlanID(d.Id()),
		talendClient.authInfo,
	)
	if err != nil {
		switch err := err.(type) {
		case *plans_executables.DeletePlanNotFound:
			fmt.Printf("%s", utils.UnmarshalErrorResponse(err.GetPayload()))
			d.SetId("") // removing from state
			return fmt.Errorf("talend plan not found, removing from state")
		}
	}

	return nil
}

func resourceTalendPlanUpdate(d *schema.ResourceData, meta interface{}) error {
	talendClient := meta.(TalendClient)

	if d.HasChange("workspace_id") ||
		d.HasChange("name") ||
		d.HasChange("steps") {

		planRequest := parsePlanRequest(d)

		updatePlanOK, err := talendClient.client.PlansExecutables.UpdatePlan(
			plans_executables.NewUpdatePlanParams().WithPlanID(d.Id()).WithBody(&planRequest),
			talendClient.authInfo,
		)
		if err != nil {
			return err
		}

		d.SetId(*updatePlanOK.GetPayload().ID)
	}

	return nil
}

func resourceTalendPlanDelete(d *schema.ResourceData, meta interface{}) error {
	talendClient := meta.(TalendClient)

	_, err := talendClient.client.PlansExecutables.DeletePlan(
		plans_executables.NewDeletePlanParams().WithPlanID(d.Id()),
		talendClient.authInfo,
	)
	if err != nil {
		switch err := err.(type) {
		case *plans_executables.DeletePlanNotFound:
			return fmt.Errorf("%s", utils.UnmarshalErrorResponse(err.GetPayload()))
		}
		return err
	}

	return nil
}

func parsePlanRequest(d *schema.ResourceData) models.PlanRequest {
	workspaceID := d.Get("workspace_id").(string)
	name := d.Get("name").(string)
	steps := []*models.Step{}
	for _, v := range d.Get("steps").([]interface{}) {
		s := v.(map[string]interface{})
		name := s["name"].(string)
		condition := s["condition"].(string)
		taskIds := []string{}
		for _, t := range s["task_ids"].([]interface{}) {
			taskIds = append(taskIds, t.(string))
		}
		steps = append(steps, &models.Step{
			Name:      &name,
			Condition: &condition,
			TaskIds:   taskIds,
		})
	}

	return models.PlanRequest{
		Name:        &name,
		WorkspaceID: &workspaceID,
		Steps:       steps,
	}

}
