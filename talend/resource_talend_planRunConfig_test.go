package talend

import (
	"fmt"
	"testing"

	"github.com/bartsimp/talend-rest-go/client/plans_executables"
	"github.com/bartsimp/talend-rest-go/utils"
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
		// fmt.Println("rs.Primary.ID=", rs.Primary.ID)
		// fmt.Println("rs.Primary.Attributes[\"plan_id\"]=", rs.Primary.Attributes["plan_id"])
		plan, err := tc.client.PlansExecutables.GetPlanRunConfiguration(
			plans_executables.NewGetPlanRunConfigurationParams().WithPlanID(rs.Primary.ID),
			tc.authInfo)
		if err != nil {
			switch err := err.(type) {
			case *plans_executables.GetPlanRunConfigurationBadRequest:
				return fmt.Errorf("%s", utils.UnmarshalErrorResponse(err.GetPayload()))
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
	environmentID := "63b56857b2b29e736cecce70" // default
	workspaceID := "63b56858b2b29e736cecce73"   // Personal
	artifactID := "63b585c56acf7f4c287d181b"
	artifactVersion := "0.1.0.20230401015724"
	taskName := sdkacctest.RandomWithPrefix("task")
	planName := sdkacctest.RandomWithPrefix("plan")
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

resource "talend_plan" "my_talend_plan_1" {
    workspace_id = %[2]q
    name         = %[6]q
    steps        {
                   name       = "step1"
                   condition  = "ALL_SUCCEEDED"
                   task_ids   = [talend_task.my_talend_task_1.id]
                 }
    steps        {
                   name       = "step2"
                   condition  = "ALL_SUCCEEDED"
                   task_ids   = [talend_task.my_talend_task_1.id]
                 }
}

resource "talend_plan_runconfig" "my_talend_plan_runconfig_1" {
    plan_id  = talend_plan.my_talend_plan_1.id
    trigger {
      type       = "ONCE"
      start_date = "2025-09-25"
      time_zone  = "Europe/London"
      at_times {
        type       = "AT_TIME"
        times      = [ "10:00" ]
        time       = "10:00"
        start_time  = "10:00"
        end_time    = "23:00"
        interval   = 10
      }
    }
    runtime {
      type = "CLOUD"
    }
}
`, environmentID, workspaceID, artifactID, artifactVersion, taskName, planName)
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
