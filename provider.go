package main

import (
	"fmt"

	"github.com/bartsimp/talend-rest-go/client"
	"github.com/go-openapi/runtime"
	"github.com/go-openapi/strfmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"api_key": {
				Type:        schema.TypeString,
				Description: "Your Talend API key",
				Required:    true,
				Sensitive:   true,
				DefaultFunc: schema.EnvDefaultFunc("TALEND_API_KEY", nil),
			},
			"host": {
				Type:        schema.TypeString,
				Description: "Host",
				// Required:    false,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("TALEND_HOST", "api.eu.cloud.talend.com"),
			},
			"schema": {
				Type:        schema.TypeString,
				Description: "Schema",
				// Required:    false,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("TALEND_SCHEMA", "https"),
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"talend_environment":    resourceTalendEnvironment(),
			"talend_workspace":      resourceTalendWorkspace(),
			"talend_task":           resourceTalendTask(),
			"talend_task_runconfig": resourceTalendTaskRunConfig(),
			"talend_plan":           resourceTalendPlan(),
			"talend_plan_runconfig": resourceTalendPlanRunConfig(),
		},
		ConfigureFunc: configureFunc(),
	}
}

func configureFunc() func(*schema.ResourceData) (interface{}, error) {
	return func(d *schema.ResourceData) (interface{}, error) {
		tc := TalendClient{
			client: client.NewHTTPClientWithConfig(
				strfmt.Default,
				client.DefaultTransportConfig().
					WithHost(d.Get("host").(string)).
					WithBasePath("/tmc/v2.6").
					WithSchemes([]string{d.Get("schema").(string)}),
			),
			authInfo: runtime.ClientAuthInfoWriterFunc(
				func(cr runtime.ClientRequest, r strfmt.Registry) error {
					cr.SetHeaderParam("Authorization", fmt.Sprint("Bearer ", d.Get("api_key").(string)))
					return nil
				},
			),
		}
		return tc, nil
	}
}

type TalendClient struct {
	client   *client.TalendManagementConsolePublicAPI
	authInfo runtime.ClientAuthInfoWriterFunc
}
