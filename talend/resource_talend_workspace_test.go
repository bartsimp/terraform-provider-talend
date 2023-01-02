package talend

import (
	"fmt"
	"testing"

	"github.com/bartsimp/talend-rest-go/client/environments"
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
	workspaceName := sdkacctest.RandomWithPrefix("ws")
	workspaceDesc := fmt.Sprintf("desc for %s", workspaceName)
	owner := "bevave5893"
	environmentID := "63a2e0dfaefa2e4ea7b1f4ae" // default
	return fmt.Sprintf(`
		resource "talend_workspace" "my_talend_workspace_1" {
			name			= %[1]q
			description		= %[2]q
			environment_id	= %[3]q
			owner			= %[4]q
		}
	`, workspaceName, workspaceDesc, environmentID, owner)
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
		if rs.Type != "talend_environment" {
			continue
		}

		query := fmt.Sprintf("name==%s", rs.Primary.Attributes["name"])
		getEnvironmentsOK, err := tc.client.Environments.GetEnvironments(
			environments.NewGetEnvironmentsParams().WithQuery(&query),
			tc.authInfo)
		if err != nil {
			return err
		}
		if getEnvironmentsOK.GetPayload() != nil {
			return fmt.Errorf("Talend Workspace still exists: %s", rs.Primary.ID)
		}
	}

	return fmt.Errorf("CheckDestroy failed")
}
