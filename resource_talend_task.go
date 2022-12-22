package main

import (
	"fmt"

	"github.com/bartsimp/talend-rest-go/client/tasks"
	"github.com/bartsimp/talend-rest-go/client/workspaces"
	"github.com/bartsimp/talend-rest-go/models"
	"github.com/bartsimp/talend-rest-go/utils"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceTalendTask() *schema.Resource {
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
			"description": {
				Type:     schema.TypeString,
				Required: true,
			},
			"artifact": {
				Type:     schema.TypeSet,
				Required: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Required: true,
						},
						"version": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
			"environment_id": {
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
	talendClient := meta.(TalendClient)

	taskRequest := parseCreate(d)

	createTaskCreated, err := talendClient.client.Tasks.CreateTask(
		tasks.NewCreateTaskParams().WithBody(&taskRequest),
		talendClient.authInfo,
	)
	if err != nil {
		switch err := err.(type) {
		case *tasks.CreateTaskBadRequest:
			return fmt.Errorf("%s", utils.UnmarshalErrorResponse(err.GetPayload()))
		}
		return err
	}

	d.SetId(*createTaskCreated.GetPayload().ID)
	return resourceTalendTaskRead(d, meta)
}

func resourceTalendTaskRead(d *schema.ResourceData, meta interface{}) error {
	talendClient := meta.(TalendClient)

	getTaskOK, err := talendClient.client.Tasks.GetTask(
		tasks.NewGetTaskParams().WithTaskID(d.Id()),
		talendClient.authInfo,
	)
	if err != nil {
		return err
	}

	d.SetId(*getTaskOK.GetPayload().ID)
	// d.Set("name", getTaskOK.GetPayload().Name)

	return nil
}

func resourceTalendTaskUpdate(d *schema.ResourceData, meta interface{}) error {
	talendClient := meta.(TalendClient)

	if d.HasChange("workspace_id") ||
		d.HasChange("name") ||
		d.HasChange("description") ||
		d.HasChange("artifact") {

		taskRequest := parseUpdate(d)

		updateTaskOK, err := talendClient.client.Tasks.UpdateTask(
			tasks.NewUpdateTaskParams().WithTaskID(d.Id()).WithBody(&taskRequest),
			talendClient.authInfo,
		)
		if err != nil {
			return err
		}

		d.SetId(*updateTaskOK.GetPayload().ID)
		// d.Set("name", updateTaskOK.GetPayload().Name)
	}

	return nil
}

func resourceTalendTaskDelete(d *schema.ResourceData, meta interface{}) error {
	talendClient := meta.(TalendClient)

	_, err := talendClient.client.Tasks.DeleteTask(
		tasks.NewDeleteTaskParams().WithTaskID(d.Id()),
		talendClient.authInfo,
	)
	if err != nil {
		switch err := err.(type) {
		case *workspaces.UpdateCustomWorkspaceBadRequest:
		case *tasks.DeleteTaskConflict:
		case *tasks.DeleteTaskNotFound:
			return fmt.Errorf("%s", utils.UnmarshalErrorResponse(err.GetPayload()))
		}
		return err
	}

	return nil
}

func parseCreate(d *schema.ResourceData) models.TaskV21CreateRequest {
	workspaceID := d.Get("workspace_id").(string)
	name := d.Get("name").(string)
	description := d.Get("description").(string)
	a := d.Get("artifact").(*schema.Set)
	aa := a.List()[0]
	artifact := aa.(map[string]interface{})
	artifactId := artifact["id"].(string)
	artifactVersion := artifact["version"].(string)
	environmentID := d.Get("environment_id").(string)

	return models.TaskV21CreateRequest{
		WorkspaceID: &workspaceID,
		Name:        &name,
		Description: &description,
		Artifact: &models.ArtifactRequest{
			ID:      &artifactId,
			Version: &artifactVersion,
		},
		EnvironmentID: &environmentID,
	}
}

func parseUpdate(d *schema.ResourceData) models.TaskAutoUpgradeRequest {
	workspaceID := d.Get("workspace_id").(string)
	name := d.Get("name").(string)
	description := d.Get("description").(string)
	a := d.Get("artifact").(*schema.Set)
	aa := a.List()[0]
	artifact := aa.(map[string]interface{})
	artifactId := artifact["id"].(string)
	artifactVersion := artifact["version"].(string)

	return models.TaskAutoUpgradeRequest{
		WorkspaceID: &workspaceID,
		Name:        &name,
		Description: &description,
		Artifact: &models.ArtifactRequest{
			ID:      &artifactId,
			Version: &artifactVersion,
		},
	}

}
