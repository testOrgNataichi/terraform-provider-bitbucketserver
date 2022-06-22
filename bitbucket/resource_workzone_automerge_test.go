package bitbucket

import (
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccBitbucketResourceWorkzoneAutomerge_requiredArgumentsOnly(t *testing.T) {
	projectKey := fmt.Sprintf("TEST%v", rand.New(rand.NewSource(time.Now().UnixNano())).Int())

	config := baseConfigForWorkzoneTests(projectKey) + `
	resource "bitbucketserver_workzone_automerge" "test" {
		project    = bitbucketserver_project.test.key
		repository = bitbucketserver_repository.test.slug
		refname    = "refs/heads/master"
	}`

	resourceName := "bitbucketserver_workzone_automerge.test"

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
					resource.TestCheckResourceAttr(resourceName, "merge_strategy_id", "none-inherit"),
					resource.TestCheckResourceAttr(resourceName, "approval_quota_enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "delete_source_branch", "false"),
					resource.TestCheckResourceAttr(resourceName, "watch_build_result", "false"),
					resource.TestCheckResourceAttr(resourceName, "watch_task_completion", "false"),
					resource.TestCheckResourceAttr(resourceName, "ignore_contributing_reviewers_approval", "true"),
					resource.TestCheckResourceAttr(resourceName, "enable_needs_work_veto", "false"),
				),
			},
		},
	})
}

func TestAccBitbucketResourceWorkzoneAutomerge_allArguments(t *testing.T) {
	projectKey := fmt.Sprintf("TEST%v", rand.New(rand.NewSource(time.Now().UnixNano())).Int())

	config := baseConfigForWorkzoneTests(projectKey) + `
	resource "bitbucketserver_workzone_automerge" "test" {
		project                                = bitbucketserver_project.test.key
		repository                             = bitbucketserver_repository.test.slug
		refname                                = "refs/heads/master"
		source_refname                         = "refs/heads/dev"
		merge_condition                        = "approvalQuota >= 5% & groupQuota >= 2"
		approval_quota                         = "6"
		group_quota                            = 3
		merge_strategy_id                      = "no-ff"
		approval_quota_enabled                 = false
		approval_count                         = 20
		mandatory_approval_count               = 10
		delete_source_branch                   = true
		watch_build_result                     = true
		watch_task_completion                  = true
		required_signatures                    = 1
		required_builds                        = 1
		ignore_contributing_reviewers_approval = false
		enable_needs_work_veto                 = true
		automerge_users                        = ["admin"]
	}`

	resourceName := "bitbucketserver_workzone_automerge.test"

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
					resource.TestCheckResourceAttr(resourceName, "source_refname", "refs/heads/dev"),
					resource.TestCheckResourceAttr(resourceName, "merge_condition", "approvalQuota >= 5% & groupQuota >= 2"),
					resource.TestCheckResourceAttr(resourceName, "merge_strategy_id", "no-ff"),
					resource.TestCheckResourceAttr(resourceName, "approval_quota_enabled", "false"),
					resource.TestCheckResourceAttr(resourceName, "approval_quota", "6"),
					resource.TestCheckResourceAttr(resourceName, "group_quota", "3"),
					resource.TestCheckResourceAttr(resourceName, "approval_count", "20"),
					resource.TestCheckResourceAttr(resourceName, "mandatory_approval_count", "10"),
					resource.TestCheckResourceAttr(resourceName, "delete_source_branch", "true"),
					resource.TestCheckResourceAttr(resourceName, "watch_build_result", "true"),
					resource.TestCheckResourceAttr(resourceName, "watch_task_completion", "true"),
					resource.TestCheckResourceAttr(resourceName, "required_signatures", "1"),
					resource.TestCheckResourceAttr(resourceName, "required_builds", "1"),
					resource.TestCheckResourceAttr(resourceName, "ignore_contributing_reviewers_approval", "false"),
					resource.TestCheckResourceAttr(resourceName, "enable_needs_work_veto", "true"),
					resource.TestCheckResourceAttr(resourceName, "automerge_users.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "automerge_users.0", "admin"),
				),
			},
		},
	})
}
