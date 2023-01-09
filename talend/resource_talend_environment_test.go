package talend

import (
	"fmt"
	"testing"

	"github.com/bartsimp/talend-rest-go/client/environments"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestTalendEnvironmentBasic(t *testing.T) {
	environmentResourceName := "talend_environment.my_talend_environment_1"
	owner := "dojon70323"
	environmentName := acctest.RandomWithPrefix("env")
	environmentDesc := fmt.Sprintf("desc for %s", environmentName)
	workspaceName := fmt.Sprintf("ws-%s", environmentName)
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testPreCheck(t) },
		ProviderFactories: testProviderFactories,
		CheckDestroy:      testTalendEnvironmentDestroy,
		Steps: []resource.TestStep{
			{
				Config: testTalendEnvironmentConfigBasic(owner, environmentName, environmentDesc, workspaceName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(environmentResourceName, "owner", owner),
					resource.TestCheckResourceAttr(environmentResourceName, "name", environmentName),
					resource.TestCheckResourceAttr(environmentResourceName, "description", environmentDesc),
					resource.TestCheckResourceAttr(environmentResourceName, "workspace_name", workspaceName),
				),
			},
		},
	})
}

func testTalendEnvironmentConfigBasic(owner, environmentName, environmentDesc, workspaceName string) string {
	return fmt.Sprintf(`
resource "talend_environment" "my_talend_environment_1" {
    owner           = %[1]q
    name            = %[2]q
    description     = %[3]q
    workspace_name  = %[4]q
}
`, owner, environmentName, environmentDesc, workspaceName)
}

func testTalendEnvironmentDestroy(s *terraform.State) error {
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
			if len(getEnvironmentsOK.GetPayload()) == 0 {
				return nil
			}
			return fmt.Errorf("Talend Environment still exists: %s", rs.Primary.ID)
		}
	}

	return fmt.Errorf("CheckDestroy failed")
}
