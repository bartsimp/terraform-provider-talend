package talend

import (
	"fmt"
	"testing"

	"github.com/bartsimp/talend-rest-go/client/plans_executables"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestTalendPlanBasic(t *testing.T) {
	taskResourceName := "talend_task.my_talend_task_1"
	planResourceName := "talend_plan.my_talend_plan_1"
	environmentID := "63b56857b2b29e736cecce70" // default
	workspaceID := "63b56858b2b29e736cecce73"   // Personal
	artifactID := "63b585c56acf7f4c287d181b"
	artifactVersion := "0.1.0.20230401015724"
	taskName := acctest.RandomWithPrefix("task")
	taskDesc := fmt.Sprintf("desc for %s", taskName)
	planName := acctest.RandomWithPrefix("plan")

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testPreCheck(t) },
		ProviderFactories: testProviderFactories,
		CheckDestroy:      testTalendPlanDestroy,
		Steps: []resource.TestStep{
			{
				Config: testTalendPlanConfigBasic(environmentID, workspaceID, artifactID, artifactVersion, taskName, taskDesc, planName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(taskResourceName, "environment_id", environmentID),
					resource.TestCheckResourceAttr(taskResourceName, "workspace_id", workspaceID),
					resource.TestCheckResourceAttr(taskResourceName, "name", taskName),
					resource.TestCheckResourceAttr(taskResourceName, "description", taskDesc),

					resource.TestCheckResourceAttr(planResourceName, "workspace_id", workspaceID),
					resource.TestCheckResourceAttr(planResourceName, "name", planName),
					resource.TestCheckResourceAttr(planResourceName, "steps.0.name", "step1"),
					resource.TestCheckResourceAttr(planResourceName, "steps.0.condition", "ALL_SUCCEEDED"),
					resource.TestCheckResourceAttr(planResourceName, "steps.1.name", "step2"),
					resource.TestCheckResourceAttr(planResourceName, "steps.1.condition", "ALL_SUCCEEDED"),
				),
			},
		},
	})
}

func testTalendPlanConfigBasic(environmentID, workspaceID, artifactID, artifactVersion, taskName, taskDesc, planName string) string {
	return fmt.Sprintf(`
resource "talend_task" "my_talend_task_1" {
    environment_id = %[1]q
    workspace_id   = %[2]q
    name           = %[5]q
    description    = %[6]q
    artifact       {
                     id      = %[3]q
                     version = %[4]q
                   }
}

resource "talend_plan" "my_talend_plan_1" {
    workspace_id = %[2]q
    name         = %[7]q
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
`, environmentID, workspaceID, artifactID, artifactVersion, taskName, taskDesc, planName)
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
