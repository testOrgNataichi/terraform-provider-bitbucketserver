package bitbucket

import (
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccBitbucketResourceWorkzoneWorkflow_requiredArgumentsOnly(t *testing.T) {
	projectKey := fmt.Sprintf("TEST%v", rand.New(rand.NewSource(time.Now().UnixNano())).Int())

	config := baseConfigForWorkzoneTests(projectKey) + `
	resource "bitbucketserver_workzone_workflow" "test" {
		project                          = bitbucketserver_project.test.key
		repository                       = bitbucketserver_repository.test.slug
	}`

	resourceName := "bitbucketserver_workzone_workflow.test"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "id", fmt.Sprintf("%v|repo", projectKey)),
					resource.TestCheckResourceAttr(resourceName, "project", projectKey),
					resource.TestCheckResourceAttr(resourceName, "repository", "repo"),
					resource.TestCheckResourceAttr(resourceName, "allow_push_after_pr", "true"),
					resource.TestCheckResourceAttr(resourceName, "unapprove_pr_after_source_change", "false"),
					resource.TestCheckResourceAttr(resourceName, "unapprove_pr_after_target_change", "false"),
					resource.TestCheckResourceAttr(resourceName, "enforce_merge_condition", "true"),
				),
			},
		},
	})
}

func TestAccBitbucketResourceWorkzoneWorkflow_allArguments(t *testing.T) {
	projectKey := fmt.Sprintf("TEST%v", rand.New(rand.NewSource(time.Now().UnixNano())).Int())

	config := baseConfigForWorkzoneTests(projectKey) + `
	resource "bitbucketserver_workzone_workflow" "test" {
		project                          = bitbucketserver_project.test.key
		repository                       = bitbucketserver_repository.test.slug
		allow_push_after_pr              = false
		unapprove_pr_after_source_change = true
		unapprove_pr_after_target_change = true
		enforce_merge_condition          = false
	}`

	resourceName := "bitbucketserver_workzone_workflow.test"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "id", fmt.Sprintf("%v|repo", projectKey)),
					resource.TestCheckResourceAttr(resourceName, "project", projectKey),
					resource.TestCheckResourceAttr(resourceName, "repository", "repo"),
					resource.TestCheckResourceAttr(resourceName, "allow_push_after_pr", "false"),
					resource.TestCheckResourceAttr(resourceName, "unapprove_pr_after_source_change", "true"),
					resource.TestCheckResourceAttr(resourceName, "unapprove_pr_after_target_change", "true"),
					resource.TestCheckResourceAttr(resourceName, "enforce_merge_condition", "false"),
				),
			},
		},
	})
}
