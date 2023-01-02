package talend

import (
	"fmt"

	"github.com/bartsimp/talend-rest-go/client/workspaces"
	"github.com/bartsimp/talend-rest-go/models"
	"github.com/bartsimp/talend-rest-go/utils"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceTalendWorkspace() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"description": {
				Type:     schema.TypeString,
				Required: true,
			},
			"environment_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"owner": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
		Create: resourceTalendWorkspaceCreate,
		Read:   resourceTalendWorkspaceRead,
		Update: resourceTalendWorkspaceUpdate,
		Delete: resourceTalendWorkspaceDelete,
	}
}

func resourceTalendWorkspaceCreate(d *schema.ResourceData, meta interface{}) error {
	talendClient := meta.(TalendClient)

	workspaceName := d.Get("name").(string)
	workspaceDesc := d.Get("description").(string)
	workspaceEnvID := d.Get("environment_id").(string)
	owner := d.Get("owner").(string)

	createCustomWorkspaceOK, createCustomWorkspaceCreated, err := talendClient.client.Workspaces.CreateCustomWorkspace(
		workspaces.NewCreateCustomWorkspaceParams().WithBody(
			&models.CreateWorkspaceRequest{
				Name:          &workspaceName,
				Description:   workspaceDesc,
				EnvironmentID: &workspaceEnvID,
				Owner:         owner,
			},
		),
		talendClient.authInfo,
	)
	if err != nil {
		switch err := err.(type) {
		case *workspaces.CreateCustomWorkspaceBadRequest:
		case *workspaces.CreateCustomWorkspaceConflict:
			return fmt.Errorf("%s", utils.UnmarshalErrorResponse(err.GetPayload()))
		}
		return err
	}

	if createCustomWorkspaceOK != nil {
		d.SetId(createCustomWorkspaceOK.GetPayload().ID)
	} else if createCustomWorkspaceCreated != nil {
		d.SetId(createCustomWorkspaceCreated.GetPayload().ID)
	}
	return nil
}

func resourceTalendWorkspaceRead(d *schema.ResourceData, meta interface{}) error {
	talendClient := meta.(TalendClient)

	query := fmt.Sprintf("name==%s;environment.id==%s", d.Get("name").(string), d.Get("environment_id").(string))
	getWorkspacesOK, err := talendClient.client.Workspaces.GetWorkspaces(
		workspaces.NewGetWorkspacesParams().WithQuery(&query),
		talendClient.authInfo,
	)
	if err != nil {
		switch err := err.(type) {
		case *workspaces.GetWorkspacesUnauthorized:
			d.SetId("") // removing from state
			return fmt.Errorf("unauthorized %s", utils.UnmarshalErrorResponse(err.GetPayload()))
		}
		return err
	}
	if len(getWorkspacesOK.GetPayload()) == 0 {
		d.SetId("") // removing from state
		return fmt.Errorf("talend workspace not found, removing from state")
	}

	return nil
}

func resourceTalendWorkspaceUpdate(d *schema.ResourceData, meta interface{}) error {
	talendClient := meta.(TalendClient)

	if d.HasChange("environment_id") ||
		d.HasChange("name") ||
		d.HasChange("description") ||
		d.HasChange("owner") {

		workspaceName := d.Get("name").(string)
		workspaceDesc := d.Get("description").(string)
		owner := d.Get("owner").(string)

		_, err := talendClient.client.Workspaces.UpdateCustomWorkspace(
			workspaces.NewUpdateCustomWorkspaceParams().WithWorkspaceID(d.Id()).WithBody(&models.UpdateWorkspaceRequest{
				Name:        &workspaceName,
				Description: workspaceDesc,
				Owner:       owner,
			}),
			talendClient.authInfo,
		)
		if err != nil {
			return err
		}
	}

	return nil
}

func resourceTalendWorkspaceDelete(d *schema.ResourceData, meta interface{}) error {
	talendClient := meta.(TalendClient)

	err := talendClient.client.Workspaces.DeleteWorkspace(
		workspaces.NewDeleteWorkspaceParams().WithWorkspaceID(d.Id()),
		talendClient.authInfo,
	)
	if err != nil {
		return err
	}

	return nil
}
