package talend

import (
	"fmt"
	"testing"

	"github.com/bartsimp/talend-rest-go/client/plans_executables"
	sdkacctest "github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestTalendPlanRunConfigBasic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testPreCheck(t) },
		Providers:    testProviders,
		CheckDestroy: testTalendPlanRunConfigDestroy,
		Steps: []resource.TestStep{
			{
				Config: testTalendPlanRunConfigConfigBasic(),
				Check: resource.ComposeTestCheckFunc(
					testTalendPlanRunConfigExists("talend_plan_runconfig.my_talend_plan_runconfig_1"),
				),
			},
		},
	})
}

func testTalendPlanRunConfigDestroy(s *terraform.State) error {
	tc := testProvider.Meta().(TalendClient)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "talend_plan_runconfig" {
			continue
		}

		plan, err := tc.client.PlansExecutables.ConfigurePlanExecution(
			plans_executables.NewConfigurePlanExecutionParams().WithPlanID(rs.Primary.ID),
			tc.authInfo)
		if err != nil {
			switch err.(type) {
			case *plans_executables.GetExecutableDetailsNotFound:
				return nil // correct, expected result
			}
			return err
		}
		if plan.GetPayload() != nil {
			return fmt.Errorf("Talend Plan RunConfig still exists for plan: %s", rs.Primary.ID)
		}
	}

	return fmt.Errorf("CheckDestroy failed")
}

func testTalendPlanRunConfigConfigBasic() string {
	// taskID := "63a4607876f4556a30a7f530"
	environmentID := "63a2e0dfaefa2e4ea7b1f4ae" // default
	workspaceID := "63a2e0dfaefa2e4ea7b1f4b1"   // Personal
	artifactID := "63a30b1d6acf7f4c287cd9e6"
	artifactVersion := "0.1.0.20222112013315"
	taskName := sdkacctest.RandomWithPrefix("task")
	return fmt.Sprintf(`
		resource "talend_task" "my_talend_task_1" {
			environment_id = %[1]q
			workspace_id   = %[2]q
			name           = %[5]q
			description    = "description for %[5]s"
			artifact       {
							id      = %[3]q
							version = %[4]q
						}
		}

		resource "talend_task_runconfig" "my_talend_task_runconfig_1" {
		  task_id  = talend_task.my_talend_task_1.id
		  trigger {
		    type       = "MANUAL"
		    start_date = "2019-09-25"
		    time_zone  = "Europe/London"
		  }
		  runtime {
		    type = "CLOUD"
		  }
		}
	`, environmentID, workspaceID, artifactID, artifactVersion, taskName)
}

func testTalendPlanRunConfigExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]

		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No TaskID set")
		}

		return nil
	}
}
