package tests

import (
	cs "cs-tf-provider/provider"
	"fmt"
	"os"
	"testing"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var testAccProvider *schema.Provider
var testAccProviders map[string]*schema.Provider
var testAccExternalProviders map[string]resource.ExternalProvider

func init() {
	testAccProvider = cs.Provider()
	testAccProviders = map[string]*schema.Provider{
		"chaossearch": testAccProvider,
	}

	testAccExternalProviders = map[string]resource.ExternalProvider{
		"aws": {
			Source: "hashicorp/aws",
		},
	}

}

func TestProvider(t *testing.T) {
	provider := cs.Provider()
	if err := provider.InternalValidate(); err != nil {
		t.Fatalf("Provider Validation Error => Error: %s", err)
	}
}

func TestProvider_impl(t *testing.T) {
	var _ schema.Provider = *cs.Provider()
}

func testAccPreCheck(t *testing.T) {
	envVars := []string{
		"CS_URL",
		"CS_ACCESS_KEY",
		"CS_SECRET_KEY",
		"CS_REGION",
		"CS_USERNAME",
		"CS_PASSWORD",
	}

	for _, val := range envVars {
		if envVar := os.Getenv(val); envVar == "" {
			t.Fatalf("Provider Configuration Error => Error: Missing %s Environment Variable", val)
		}
	}
}

func testAccProviderConfigBlock() string {
	return `
		provider "chaossearch" {
		  login {}
	    }
	`
}

func generateName(bucketName string) string {
	return fmt.Sprintf("%s-%s", bucketName, uuid.New())
}
