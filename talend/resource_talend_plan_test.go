package talend

import (
	"fmt"
	"testing"

	"github.com/bartsimp/talend-rest-go/client/plans_executables"
	sdkacctest "github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestTalendPlanBasic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testPreCheck(t) },
		Providers:    testProviders,
		CheckDestroy: testTalendPlanDestroy,
		Steps: []resource.TestStep{
			{
				Config: testTalendPlanConfigBasic(),
				Check: resource.ComposeTestCheckFunc(
					testTalendPlanExists("talend_plan.my_talend_plan_1"),
				),
			},
		},
	})
}

func testTalendPlanDestroy(s *terraform.State) error {
	tc := testProvider.Meta().(TalendClient)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "talend_plan" {
			continue
		}

		plan, err := tc.client.PlansExecutables.GetExecutableDetails(
			plans_executables.NewGetExecutableDetailsParams().WithPlanID(rs.Primary.ID),
			tc.authInfo,
		)
		if err != nil {
			switch err.(type) {
			case *plans_executables.GetExecutableDetailsNotFound:
				return nil // correct, expected result
			}
			return err
		}
		if plan.GetPayload() != nil {
			return fmt.Errorf("Talend Plan still exists: %s", rs.Primary.ID)
		}
	}
	return fmt.Errorf("CheckDestroy failed")
}

func testTalendPlanConfigBasic() string {
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
`, environmentID, workspaceID, artifactID, artifactVersion, taskName, planName)
}

func testTalendPlanExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]

		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No PlanID set")
		}

		return nil
	}
}
