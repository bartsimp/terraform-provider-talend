package talend

import (
	"fmt"
	"testing"

	"github.com/bartsimp/talend-rest-go/client/tasks"
	sdkacctest "github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestTalendTaskBasic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testPreCheck(t) },
		Providers:    testProviders,
		CheckDestroy: testTalendTaskDestroy,
		Steps: []resource.TestStep{
			{
				Config: testTalendTaskConfigBasic(),
				Check: resource.ComposeTestCheckFunc(
					testTalendTaskExists("talend_task.my_talend_task_1"),
				),
			},
		},
	})
}

func testTalendTaskDestroy(s *terraform.State) error {
	tc := testProvider.Meta().(TalendClient)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "talend_task" {
			continue
		}

		task, err := tc.client.Tasks.GetTask(
			tasks.NewGetTaskParams().WithTaskID(rs.Primary.ID),
			tc.authInfo)
		if err != nil {
			switch err.(type) {
			case *tasks.GetTaskNotFound:
				return nil // correct, expected result
			}
			return err
		}
		if task.GetPayload() != nil {
			return fmt.Errorf("Talend Task still exists: %s", rs.Primary.ID)
		}
	}

	return fmt.Errorf("CheckDestroy failed")
}

func testTalendTaskConfigBasic() string {
	environmentID := "63a2e0dfaefa2e4ea7b1f4ae" // default
	workspaceID := "63a2e0dfaefa2e4ea7b1f4b1"   // Personal
	artifactID := "63a30b1d6acf7f4c287cd9e6"
	artifactVersion := "0.1.0.20222112013315"
	taskName := sdkacctest.RandomWithPrefix("task")
	return fmt.Sprintf(`
		resource "talend_task" "my_talend_task_1" {
			environment_id	= %[1]q
			workspace_id	= %[2]q
			name			= %[5]q
			description		= "description for %[5]s"
			artifact		{
				id		= %[3]q
				version	= %[4]q
			}
			
		}
	`, environmentID, workspaceID, artifactID, artifactVersion, taskName)
}

func testTalendTaskExists(n string) resource.TestCheckFunc {
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