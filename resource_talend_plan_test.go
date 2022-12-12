package main

import (
	"fmt"
	"testing"

	"github.com/bartsimp/talend-rest-go/client/plans_executables"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestTalendPlanBasic(t *testing.T) {
	workspaceID := "myWsID"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testPreCheck(t) },
		Providers:    testProviders,
		CheckDestroy: testTalendPlanDestroy,
		Steps: []resource.TestStep{
			{
				Config: testTalendPlanConfigBasic(workspaceID),
				Check: resource.ComposeTestCheckFunc(
					testTalendPlanExists("talend_plan.my_talend_plan_1"),
				),
			},
		},
	})
}

func testTalendPlanDestroy(s *terraform.State) error {
	tc := testProvider.Meta().(TalendClient)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "talend_plan" {
			continue
		}

		planID := rs.Primary.ID

		_, err := tc.client.PlansExecutables.DeletePlan(
			plans_executables.NewDeletePlanParams().WithPlanID(planID),
			tc.authInfo,
		)
		if err != nil {
			return err
		}
	}
	return nil
}

func testTalendPlanConfigBasic(workspaceID string) string {
	return fmt.Sprintf(`
		resource "talend_plan" "my_talend_plan_1" {
			workspace_id	= "%s"
			name            = "simple executable"
			steps {
					name       = "step1"
					condition  = "ALL_SUCCEEDED"
					task_ids   = ["57f64991e4b0b689a64feed3", "57f64991e4b0b689a64feed4"]
			}
			steps {
					name       = "step2"
					condition  = "ALL_SUCCEEDED"
					task_ids   = ["57f64991e4b0b689a64feed5", "57f64991e4b0b689a64feed6"]
			}
		}
	`, workspaceID)
}

func testTalendPlanExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]

		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No PlanID set")
		}

		return nil
	}
}
