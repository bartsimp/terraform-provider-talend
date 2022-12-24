package main

import (
	"fmt"

	"github.com/bartsimp/talend-rest-go/client/environments"
	"github.com/bartsimp/talend-rest-go/models"
	"github.com/bartsimp/talend-rest-go/utils"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceTalendEnvironment() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"workspace_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"owner": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
		Create: resourceTalendEnvironmentCreate,
		Read:   resourceTalendEnvironmentRead,
		Update: resourceTalendEnvironmentUpdate,
		Delete: resourceTalendEnvironmentDelete,
	}
}

func resourceTalendEnvironmentCreate(d *schema.ResourceData, meta interface{}) error {
	talendClient := meta.(TalendClient)

	environmentName := d.Get("name").(string)
	environmentDesc := d.Get("description").(string)
	workspaceName := d.Get("workspace_name").(string)
	owner := d.Get("owner").(string)

	createCustomEnvironmentCreated, err := talendClient.client.Environments.CreateEnvironment(
		environments.NewCreateEnvironmentParams().WithBody(
			&models.CreateEnvironmentRequest{
				Name:          &environmentName,
				Description:   environmentDesc,
				WorkspaceName: &workspaceName,
				Owner:         &owner,
			},
		),
		talendClient.authInfo,
	)
	if err != nil {
		switch err := err.(type) {
		case *environments.CreateEnvironmentBadRequest:
		case *environments.CreateEnvironmentConflict:
			return fmt.Errorf("%s", utils.UnmarshalErrorResponse(err.GetPayload()))
		}
		return err
	}

	d.SetId(*createCustomEnvironmentCreated.GetPayload().ID)
	return nil
}

func resourceTalendEnvironmentRead(d *schema.ResourceData, meta interface{}) error {
	talendClient := meta.(TalendClient)

	query := fmt.Sprintf("name==%s", d.Get("name").(string))
	getEnvironmentsOK, err := talendClient.client.Environments.GetEnvironments(
		environments.NewGetEnvironmentsParams().WithQuery(&query),
		talendClient.authInfo,
	)
	if err != nil {
		switch err := err.(type) {
		case *environments.GetEnvironmentsUnauthorized:
			d.SetId("") // removing from state
			return fmt.Errorf("unauthorized %s", utils.UnmarshalErrorResponse(err.GetPayload()))
		}
		return err
	}
	if len(getEnvironmentsOK.GetPayload()) == 0 {
		d.SetId("") // removing from state
		return fmt.Errorf("talend environment not found, removing from state")
	}

	return nil
}

func resourceTalendEnvironmentUpdate(d *schema.ResourceData, meta interface{}) error {
	talendClient := meta.(TalendClient)

	if d.HasChange("name") ||
		d.HasChange("description") {

		environmentName := d.Get("name").(string)
		environmentDesc := d.Get("description").(string)

		_, err := talendClient.client.Environments.UpdateEnvironment(
			environments.NewUpdateEnvironmentParams().WithEnvironmentID(d.Id()).WithBody(&models.UpdateEnvironmentRequest{
				Name:        &environmentName,
				Description: environmentDesc,
			}),
			talendClient.authInfo,
		)
		if err != nil {
			return err
		}
	}

	return nil
}

func resourceTalendEnvironmentDelete(d *schema.ResourceData, meta interface{}) error {
	talendClient := meta.(TalendClient)

	_, err := talendClient.client.Environments.DeleteEnvironment(
		environments.NewDeleteEnvironmentParams().WithEnvironmentID(d.Id()),
		talendClient.authInfo,
	)
	if err != nil {
		return err
	}

	return nil
}
