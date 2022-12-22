package main

import (
	"fmt"
	"testing"

	"github.com/bartsimp/talend-rest-go/client/tasks"
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

		taskID := rs.Primary.ID

		task, err := tc.client.Tasks.GetTask(
			tasks.NewGetTaskParams().WithTaskID(taskID),
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
	environmentID := "63a2e0dfaefa2e4ea7b1f4ae"
	workspaceID := "63a2e0dfaefa2e4ea7b1f4b1"
	artifactID := "63a30b1d6acf7f4c287cd9e6"
	artifactVersion := "0.1.0.20222112013315"
	return fmt.Sprintf(`
		resource "talend_task" "my_talend_task_1" {
			environment_id	= "%s"
			workspace_id	= "%s"
			name			= "Hello world task"
			description		= "Task detail description"
			artifact		{
				id		= "%s"
				version	= "%s"
			}
			
		}
	`, environmentID, workspaceID, artifactID, artifactVersion)
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
