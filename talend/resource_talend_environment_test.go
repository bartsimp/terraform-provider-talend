package talend

import (
	"fmt"
	"testing"

	"github.com/bartsimp/talend-rest-go/client/environments"
	sdkacctest "github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestTalendEnvironmentBasic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testPreCheck(t) },
		Providers:    testProviders,
		CheckDestroy: testTalendEnvironmentDestroy,
		Steps: []resource.TestStep{
			{
				Config: testTalendEnvironmentConfigBasic(),
				Check: resource.ComposeTestCheckFunc(
					testTalendEnvironmentExists("talend_environment.my_talend_environment_1"),
				),
			},
		},
	})
}

func testTalendEnvironmentConfigBasic() string {
	environmentName := sdkacctest.RandomWithPrefix("env")
	environmentDesc := fmt.Sprintf("desc for %s", environmentName)
	workspaceName := fmt.Sprintf("ws-%s", environmentName)
	owner := "dojon70323"
	return fmt.Sprintf(`
resource "talend_environment" "my_talend_environment_1" {
    name            = %[1]q
    description     = %[2]q
    workspace_name  = %[3]q
    owner           = %[4]q
}
`, environmentName, environmentDesc, workspaceName, owner)
}

func testTalendEnvironmentExists(n string) resource.TestCheckFunc {
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
