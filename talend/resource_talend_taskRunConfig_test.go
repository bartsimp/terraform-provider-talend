package talend

import (
	"fmt"
	"testing"

	"github.com/bartsimp/talend-rest-go/client/tasks"
	"github.com/bartsimp/talend-rest-go/utils"
	"github.com/go-openapi/runtime"
	sdkacctest "github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestTalendTaskRunConfigBasic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testPreCheck(t) },
		Providers:    testProviders,
		CheckDestroy: testTalendTaskRunConfigDestroy,
		Steps: []resource.TestStep{
			{
				Config: testTalendTaskRunConfigConfigBasic(),
				Check: resource.ComposeTestCheckFunc(
					testTalendTaskRunConfigExists("talend_task_runconfig.my_talend_task_runconfig_1"),
				),
			},
		},
	})
}

func testTalendTaskRunConfigDestroy(s *terraform.State) error {
	tc := testProvider.Meta().(TalendClient)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "talend_task_runconfig" {
			continue
		}
		// fmt.Println("rs.Primary.ID=", rs.Primary.ID)
		// fmt.Println("rs.Primary.Attributes[\"task_id\"]=", rs.Primary.Attributes["task_id"])
		_, err := tc.client.Tasks.GetTaskConfiguration(
			tasks.NewGetTaskConfigurationParams().WithTaskID(rs.Primary.ID),
			func(co *runtime.ClientOperation) {
				co.AuthInfo = tc.authInfo
			},
		)
		if err != nil {
			switch err := err.(type) {
			case *tasks.GetTaskConfigurationUnauthorized:
				return fmt.Errorf("unauthorized %s", utils.UnmarshalErrorResponse(err.GetPayload()))
			case *tasks.GetTaskConfigurationNotFound:
				return fmt.Errorf("%s", utils.UnmarshalErrorResponse(err.GetPayload()))
			case *tasks.GetTaskNotFound:
				return nil // correct, expected result
			}
			return err
		}
	}

	return fmt.Errorf("CheckDestroy failed")
}

func testTalendTaskRunConfigConfigBasic() string {
	environmentID := "63b56857b2b29e736cecce70" // default
	workspaceID := "63b56858b2b29e736cecce73"   // Personal
	artifactID := "63b585c56acf7f4c287d181b"
	artifactVersion := "0.1.0.20230401015724"
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
`, environmentID, workspaceID, artifactID, artifactVersion, taskName)
}

func testTalendTaskRunConfigExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		/*
			_, err := talendClient.client.Tasks.GetTaskConfiguration(
				tasks.NewGetTaskConfigurationParams().WithTaskID(d.Id()),
				func(co *runtime.ClientOperation) {
					co.AuthInfo = talendClient.authInfo
				},
			)
		*/
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No TaskID set")
		}

		return nil
	}
}
