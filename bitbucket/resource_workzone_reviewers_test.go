package bitbucket

import (
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccBitbucketResourceWorkzoneReviewers_requiredArgumentsOnly(t *testing.T) {
	projectKey := fmt.Sprintf("TEST%v", rand.New(rand.NewSource(time.Now().UnixNano())).Int())

	config := baseConfigForWorkzoneTests(projectKey) + `
	resource "bitbucketserver_workzone_reviewers" "test" {
		project    = bitbucketserver_project.test.key
		repository = bitbucketserver_repository.test.slug
		refname    = "refs/heads/master"
	}`

	resourceName := "bitbucketserver_workzone_reviewers.test"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "id", fmt.Sprintf("%v|repo|refs/heads/master", projectKey)),
					resource.TestCheckResourceAttr(resourceName, "project", projectKey),
					resource.TestCheckResourceAttr(resourceName, "repository", "repo"),
					resource.TestCheckResourceAttr(resourceName, "refname", "refs/heads/master"),
					resource.TestCheckResourceAttr(resourceName, "suggested_reviewers", "0"),
					resource.TestCheckResourceAttr(resourceName, "suggested_reviewers_timespan", "90"),
				),
			},
		},
	})
}

func TestAccBitbucketResourceWorkzoneReviewers_allArguments(t *testing.T) {
	projectKey := fmt.Sprintf("TEST%v", rand.New(rand.NewSource(time.Now().UnixNano())).Int())

	config := baseConfigForWorkzoneTests(projectKey) + `
	resource "bitbucketserver_group" "test" {
		name = "test_group"
	}

	resource "bitbucketserver_group" "test2" {
		name = "test_group_2"
	}

	resource "bitbucketserver_workzone_reviewers" "test" {
		project                               = bitbucketserver_project.test.key
		repository                            = bitbucketserver_repository.test.slug
		refname                               = "refs/heads/master"
		repository_reviewers_users            = ["admin"]
		repository_reviewers_groups           = [bitbucketserver_group.test.name]
		mandatory_repository_reviewers_users  = ["admin"]
		mandatory_repository_reviewers_groups = [bitbucketserver_group.test.name]
		suggested_reviewers                   = 2
		suggested_reviewers_timespan          = 100
		filepath_reviewers {
			filepath_pattern                    = "**/11"
			filepath_reviewers_users            = ["admin"]
			filepath_reviewers_groups           = [bitbucketserver_group.test.name]
			mandatory_filepath_reviewers_users  = ["admin"]
			mandatory_filepath_reviewers_groups = [bitbucketserver_group.test.name]
			include_exclude_enabled             = true
			exclude_filepaths                   = ["**/11/*.tfstate"]
		}
		filepath_reviewers {
			filepath_pattern                    = "**/22"
			filepath_reviewers_users            = ["admin"]
			filepath_reviewers_groups           = [bitbucketserver_group.test.name, bitbucketserver_group.test2.name]
			mandatory_filepath_reviewers_users  = ["admin"]
			mandatory_filepath_reviewers_groups = [bitbucketserver_group.test.name, bitbucketserver_group.test2.name]
			include_exclude_enabled             = false
			exclude_filepaths                   = ["**/22/*.tfstate", "**/22/*.tfstate2"]
		}
	}`

	resourceName := "bitbucketserver_workzone_reviewers.test"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "id", fmt.Sprintf("%v|repo|refs/heads/master", projectKey)),
					resource.TestCheckResourceAttr(resourceName, "project", projectKey),
					resource.TestCheckResourceAttr(resourceName, "repository", "repo"),
					resource.TestCheckResourceAttr(resourceName, "refname", "refs/heads/master"),
					resource.TestCheckResourceAttr(resourceName, "suggested_reviewers", "2"),
					resource.TestCheckResourceAttr(resourceName, "suggested_reviewers_timespan", "100"),
					resource.TestCheckResourceAttr(resourceName, "repository_reviewers_users.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "repository_reviewers_users.0", "admin"),
					resource.TestCheckResourceAttr(resourceName, "repository_reviewers_groups.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "repository_reviewers_groups.0", "test_group"),
					resource.TestCheckResourceAttr(resourceName, "mandatory_repository_reviewers_users.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "mandatory_repository_reviewers_users.0", "admin"),
					resource.TestCheckResourceAttr(resourceName, "mandatory_repository_reviewers_groups.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "mandatory_repository_reviewers_groups.0", "test_group"),

					resource.TestCheckResourceAttr(resourceName, "filepath_reviewers.#", "2"),

					resource.TestCheckResourceAttr(resourceName, "filepath_reviewers.0.filepath_pattern", "**/11"),
					resource.TestCheckResourceAttr(resourceName, "filepath_reviewers.0.filepath_reviewers_users.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "filepath_reviewers.0.filepath_reviewers_users.0", "admin"),
					resource.TestCheckResourceAttr(resourceName, "filepath_reviewers.0.filepath_reviewers_groups.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "filepath_reviewers.0.filepath_reviewers_groups.0", "test_group"),
					resource.TestCheckResourceAttr(resourceName, "filepath_reviewers.0.mandatory_filepath_reviewers_users.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "filepath_reviewers.0.mandatory_filepath_reviewers_users.0", "admin"),
					resource.TestCheckResourceAttr(resourceName, "filepath_reviewers.0.mandatory_filepath_reviewers_groups.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "filepath_reviewers.0.mandatory_filepath_reviewers_groups.0", "test_group"),
					resource.TestCheckResourceAttr(resourceName, "filepath_reviewers.0.include_exclude_enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "filepath_reviewers.0.exclude_filepaths.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "filepath_reviewers.0.exclude_filepaths.0", "**/11/*.tfstate"),

					resource.TestCheckResourceAttr(resourceName, "filepath_reviewers.1.filepath_pattern", "**/22"),
					resource.TestCheckResourceAttr(resourceName, "filepath_reviewers.1.filepath_reviewers_users.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "filepath_reviewers.1.filepath_reviewers_users.0", "admin"),
					resource.TestCheckResourceAttr(resourceName, "filepath_reviewers.1.filepath_reviewers_groups.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "filepath_reviewers.1.filepath_reviewers_groups.0", "test_group"),
					resource.TestCheckResourceAttr(resourceName, "filepath_reviewers.1.filepath_reviewers_groups.1", "test_group_2"),
					resource.TestCheckResourceAttr(resourceName, "filepath_reviewers.1.mandatory_filepath_reviewers_users.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "filepath_reviewers.1.mandatory_filepath_reviewers_users.0", "admin"),
					resource.TestCheckResourceAttr(resourceName, "filepath_reviewers.1.mandatory_filepath_reviewers_groups.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "filepath_reviewers.1.mandatory_filepath_reviewers_groups.0", "test_group"),
					resource.TestCheckResourceAttr(resourceName, "filepath_reviewers.1.mandatory_filepath_reviewers_groups.1", "test_group_2"),
					resource.TestCheckResourceAttr(resourceName, "filepath_reviewers.1.include_exclude_enabled", "false"),
					resource.TestCheckResourceAttr(resourceName, "filepath_reviewers.1.exclude_filepaths.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "filepath_reviewers.1.exclude_filepaths.0", "**/22/*.tfstate"),
					resource.TestCheckResourceAttr(resourceName, "filepath_reviewers.1.exclude_filepaths.1", "**/22/*.tfstate2"),
				),
			},
		},
	})
}
