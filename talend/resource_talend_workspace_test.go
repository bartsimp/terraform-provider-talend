package talend

import (
	"fmt"
	"testing"

	"github.com/bartsimp/talend-rest-go/client/workspaces"
	sdkacctest "github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestTalendWorkspaceBasic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testPreCheck(t) },
		Providers:    testProviders,
		CheckDestroy: testTalendWorkspaceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testTalendWorkspaceConfigBasic(),
				Check: resource.ComposeTestCheckFunc(
					testTalendWorkspaceExists("talend_workspace.my_talend_workspace_1"),
				),
			},
		},
	})
}

func testTalendWorkspaceConfigBasic() string {
	environmentName := sdkacctest.RandomWithPrefix("env")
	environmentDesc := fmt.Sprintf("desc for %s", environmentName)
	environmentWorkspaceName := fmt.Sprintf("ws-%s", environmentName)
	workspaceName := sdkacctest.RandomWithPrefix("ws")
	workspaceDesc := fmt.Sprintf("desc for %s", workspaceName)
	owner := "dojon70323"
	return fmt.Sprintf(`
resource "talend_environment" "my_talend_environment_1" {
  name            = %[1]q
  description     = %[2]q
  workspace_name  = %[3]q
  owner           = %[6]q
}

resource "talend_workspace" "my_talend_workspace_1" {
  name            = %[4]q
  description     = %[5]q
  environment_id  = talend_environment.my_talend_environment_1.id
  owner           = %[6]q
}
`, environmentName, environmentDesc, environmentWorkspaceName, workspaceName, workspaceDesc, owner)
}

func testTalendWorkspaceExists(n string) resource.TestCheckFunc {
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

func testTalendWorkspaceDestroy(s *terraform.State) error {
	tc := testProvider.Meta().(TalendClient)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "talend_workspace" {
			continue
		}

		query := fmt.Sprintf("name==%s", rs.Primary.Attributes["name"])
		getWorkspacesOK, err := tc.client.Workspaces.GetWorkspaces(
			workspaces.NewGetWorkspacesParams().WithQuery(&query),
			tc.authInfo)
		if err != nil {
			return err
		}
		if getWorkspacesOK.GetPayload() != nil {
			if len(getWorkspacesOK.GetPayload()) == 0 {
				return nil
			}
			return fmt.Errorf("Talend Workspace still exists: %s", rs.Primary.ID)
		}
	}

	return fmt.Errorf("CheckDestroy failed")
}
