package talend

import (
	talendRest "github.com/bartsimp/talend-rest-go"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
)

func Provider() terraform.ResourceProvider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"api_key": &schema.Schema{
				Type:        schema.TypeString,
				Description: "Your Talend API key",
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("TALEND_API_KEY", nil),
			},
			"proxy": &schema.Schema{
				Type:        schema.TypeString,
				Description: "Proxy url",
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("TALEND_PROXY", nil),
			},
		},
		ResourcesMap: map[string]*schema.Resource{
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
		client, err := talendRest.NewClient(d.Get("api_key").(string), d.Get("proxy").(string))
		return client, err
	}
}
