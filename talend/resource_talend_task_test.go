package talend

import (
	"fmt"
	"testing"

	"github.com/bartsimp/talend-rest-go/client/tasks"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestTalendTaskBasic(t *testing.T) {
	taskResourceName := "talend_task.my_talend_task_1"
	environmentID := "63b56857b2b29e736cecce70" // default
	workspaceID := "63b56858b2b29e736cecce73"   // Personal
	artifactID := "63b585c56acf7f4c287d181b"
	artifactVersion := "0.1.0.20230401015724"
	taskName := acctest.RandomWithPrefix("task")
	taskDesc := fmt.Sprintf("desc for %s", taskName)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testPreCheck(t) },
		Providers:    testProviders,
		CheckDestroy: testTalendTaskDestroy,
		Steps: []resource.TestStep{
			{
				Config: testTalendTaskConfigBasic(environmentID, workspaceID, artifactID, artifactVersion, taskName, taskDesc),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(taskResourceName, "environment_id", environmentID),
					resource.TestCheckResourceAttr(taskResourceName, "workspace_id", workspaceID),
					resource.TestCheckResourceAttr(taskResourceName, "name", taskName),
					resource.TestCheckResourceAttr(taskResourceName, "description", taskDesc),
					resource.TestCheckResourceAttr(taskResourceName, "artifact.0.id", artifactID),
					resource.TestCheckResourceAttr(taskResourceName, "artifact.0.version", artifactVersion),
				),
			},
		},
	})
}

func testTalendTaskConfigBasic(environmentID, workspaceID, artifactID, artifactVersion, taskName, taskDescription string) string {
	return fmt.Sprintf(`
		resource "talend_task" "my_talend_task_1" {
			environment_id	= %[1]q
			workspace_id	= %[2]q
			name			= %[5]q
			description		= %[6]q
			artifact		{
				id		= %[3]q
				version	= %[4]q
			}
			
		}
	`, environmentID, workspaceID, artifactID, artifactVersion, taskName, taskDescription)
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
