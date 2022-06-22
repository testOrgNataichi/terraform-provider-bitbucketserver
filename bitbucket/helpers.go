package bitbucket

import (
	"fmt"

	"github.com/hashicorp/terraform/helper/schema"
)

func stringArrayFromSchemaSet(schemaSet *schema.Set) []string {
	array := make([]string, 0, len(schemaSet.List()))

	for _, item := range schemaSet.List() {
		array = append(array, item.(string))
	}

	return array
}

func baseConfigForRepositoryBasedTests(projectKey string) string {
	config := fmt.Sprintf(`
		resource "bitbucketserver_project" "test" {
			key = "%v"
			name = "test-project-%v"
		}

		resource "bitbucketserver_repository" "test" {
			project = bitbucketserver_project.test.key
			name = "repo"
		}
	`, projectKey, projectKey)

	return config
}

func baseConfigForWorkzoneTests(projectKey string) string {
	config := fmt.Sprintf(`
	resource "bitbucketserver_plugin" "test" {
		key     = "com.izymes.workzone"
		version = "7.5.3"
		enabled = true
		# Get license here https://developer.atlassian.com/platform/marketplace/timebomb-licenses-for-testing-server-apps/
		license = "AAABCA0ODAoPeNpdj01PwkAURffzKyZxZ1IyUzARkllQ24gRaQMtGnaP8VEmtjPNfFT59yJVFyzfubkn796Ux0Bz6SmbUM5nbDzj97RISxozHpMUnbSq88poUaLztFEStUN6MJZ2TaiVpu/YY2M6tI6sQrtHmx8qd74EZ+TBIvyUU/AoYs7jiE0jzknWQxMuifA2IBlUbnQ7AulVjwN9AaU9atASs69O2dNFU4wXJLc1aOUGw9w34JwCTTZoe7RPqUgep2X0Vm0n0fNut4gSxl/Jcnj9nFb6Q5tP/Ueu3L+0PHW4ghZFmm2zZV5k6/95CbR7Y9bYGo/zGrV3Ir4jRbDyCA6vt34DO8p3SDAsAhQnJjLD5k9Fr3uaIzkXKf83o5vDdQIUe4XequNCC3D+9ht9ZYhNZFKmnhc=X02dh"
	}

	resource "bitbucketserver_project" "test" {
		key  = "%v"
		name = "test-project-%v"
		depends_on = [bitbucketserver_plugin.test]
	}

	resource "bitbucketserver_repository" "test" {
		project    = bitbucketserver_project.test.key
		name       = "repo"
	}`, projectKey, projectKey)
	return config
}
