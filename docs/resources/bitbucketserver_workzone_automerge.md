# Resource: bitbucketserver_workzone_automerge

Provides the ability to manage Workzone auto-merge settings. *If you omit the `automerge_user`, the settings can be used as review quotas for the pull requests.*

## Example Usage

```hcl
resource "bitbucketserver_workzone_automerge" "test" {
  project                                = "MYPROJ"
  repository                             = "repo"
  refname                                = "refs/heads/master"
  source_refname                         = "refs/heads/dev"
  merge_condition                        = "approvalQuota >= 5% & groupQuota >= 2"
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
}
```

## Argument Reference

* `project` - Required. Project Key that contains target repository.
* `repository` - Required. Repository slug of target repository.
* `refname` - Required. Refname of the destination branch for the pull requests.
* `source_refname` - Optional. Refname of the source branch for the pull requests.
* `merge_condition` - Optional. [Boolean expression-based merge condition](https://www.izymes.com/products/workzone/merge-control/).
* `merge_strategy_id` - Optional. Merge strategy. Git merge strategies affect the way the Git history appears after merging a pull request. Default `none-inherit`.
* `approval_quota_enabled` - Optional. The setting allows enabling/disabling the approval quota for the pull requests to the target branch.
* `approval_quota` - Optional. The percentage of approvals among all reviewers assigned to a pull request. It's preferable to set the value in `merge_condition` expression instead.
* `approval_count` - Optional. An absolute number of reviewer approvals required to merge the PR. It's preferable to set the value in `merge_condition` expression instead.
* `mandatory_approval_count` - Optional. The absolute number of mandatory reviewer approvals required to merge the pull request.
* `group_quota` - Optional. The number of members per reviewer group required to merge the pull request.
* `delete_source_branch` - Optional. Delete source branch after successful auto-merge. Default `false`.
* `watch_build_result` - Optional. A successful build triggers a merge if all other merge conditions are met. Default `false`.
* `required_signatures` - Optional. The number of digital signature approvals required to merge the pull request.
* `required_builds` - Optional. The number of successful builds required to merge the pull request.
* `ignore_contributing_reviewers_approval` - Optional. The setting controls if pull request code contributor approvals are counted or not towards quotas. Default `true`.
* `enable_needs_work_veto` - Optional. Whether or not a pull request with 'needs work' flag will be blocked from merging. Default `false`.
* `automerge_users` - Optional. Merge pull request as this user. The user must have write permissions to the target branch. *If you omit the `automerge_user`, the settings can be used as review quotas for the pull requests.*
