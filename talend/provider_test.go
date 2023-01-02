package talend

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var testProviders map[string]*schema.Provider
var testProvider *schema.Provider

func init() {
	testProvider = Provider()
	testProviders = map[string]*schema.Provider{
		"talend": testProvider,
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
