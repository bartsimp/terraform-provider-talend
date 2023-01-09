package talend

import (
	"fmt"
	"testing"

	"github.com/bartsimp/talend-rest-go/client/workspaces"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestTalendWorkspaceBasic(t *testing.T) {
	environmentResourceName := "talend_environment.my_talend_environment_1"
	workspaceResourceName := "talend_workspace.my_talend_workspace_1"
	owner := "dojon70323"
	environmentName := acctest.RandomWithPrefix("env")
	environmentDesc := fmt.Sprintf("desc for %s", environmentName)
	environmentWorkspaceName := fmt.Sprintf("ws-%s", environmentName)
	workspaceName := acctest.RandomWithPrefix("ws")
	workspaceDesc := fmt.Sprintf("desc for %s", workspaceName)
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testPreCheck(t) },
		ProviderFactories: testProviderFactories,
		CheckDestroy:      testTalendWorkspaceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testTalendWorkspaceConfigBasic(owner, environmentName, environmentDesc, environmentWorkspaceName, workspaceName, workspaceDesc),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(environmentResourceName, "owner", owner),
					resource.TestCheckResourceAttr(environmentResourceName, "name", environmentName),
					resource.TestCheckResourceAttr(environmentResourceName, "description", environmentDesc),
					resource.TestCheckResourceAttr(environmentResourceName, "workspace_name", environmentWorkspaceName),

					resource.TestCheckResourceAttr(workspaceResourceName, "owner", owner),
					resource.TestCheckResourceAttr(workspaceResourceName, "name", workspaceName),
					resource.TestCheckResourceAttr(workspaceResourceName, "description", workspaceDesc),
				),
			},
		},
	})
}

func testTalendWorkspaceConfigBasic(owner, environmentName, environmentDesc, environmentWorkspaceName, workspaceName, workspaceDesc string) string {
	return fmt.Sprintf(`
resource "talend_environment" "my_talend_environment_1" {
  owner           = %[1]q
  name            = %[2]q
  description     = %[3]q
  workspace_name  = %[4]q
}

resource "talend_workspace" "my_talend_workspace_1" {
  owner           = %[1]q
  name            = %[5]q
  description     = %[6]q
  environment_id  = talend_environment.my_talend_environment_1.id
}
`, owner, environmentName, environmentDesc, environmentWorkspaceName, workspaceName, workspaceDesc)
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
