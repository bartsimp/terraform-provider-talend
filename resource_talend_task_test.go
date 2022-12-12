package main

import (
	"fmt"
	"testing"

	"github.com/bartsimp/talend-rest-go/client/tasks"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestTalendTaskBasic(t *testing.T) {
	workspaceID := "myWsID"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testPreCheck(t) },
		Providers:    testProviders,
		CheckDestroy: testTalendTaskDestroy,
		Steps: []resource.TestStep{
			{
				Config: testTalendTaskConfigBasic(workspaceID),
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

		_, err := tc.client.Tasks.DeleteTask(
			tasks.NewDeleteTaskParams().WithTaskID(taskID),
			tc.authInfo,
		)
		if err != nil {
			return err
		}
	}

	return nil
}

func testTalendTaskConfigBasic(workspaceID string) string {
	return fmt.Sprintf(`
		resource "talend_task" "my_talend_task_1" {
			workspace_id	= "%s"
			name			= "Hello world task"
			description		= "Task detail description"
			artifact		{
				id		= "5c1111d7a4186a4eafed0587"
				version	= "0.1.0"
			}
			environment_id	= "5d7a3d082d909b386943787e"
		}
	`, workspaceID)
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
