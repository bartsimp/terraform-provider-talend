package talend

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var testProvider *schema.Provider
var testProviderFactories map[string]func() (*schema.Provider, error)

func init() {
	testProvider = Provider()
	testProviderFactories = map[string]func() (*schema.Provider, error){
		"talend": func() (*schema.Provider, error) {
			return testProvider, nil
		},
	}
}

func TestProvider(t *testing.T) {
	if err := Provider().InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}

func TestProvider_impl(t *testing.T) {
	var _ *schema.Provider = Provider()
}

func testPreCheck(t *testing.T) {
	if err := os.Getenv("TALEND_API_KEY"); err == "" {
		t.Fatal("TALEND_API_KEY must be set for acceptance tests")
	}
}
