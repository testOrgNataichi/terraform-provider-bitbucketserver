# Resource: bitbucketserver_workzone_workflow

Provides the ability to manage Workzone workflow settings.

## Example Usage

```hcl
resource "bitbucketserver_workzone_workflow" "test" {
  project                          = "MYPROJ"
  repository                       = "repo"
  allow_push_after_pr              = false
  unapprove_pr_after_source_change = true
  unapprove_pr_after_target_change = true
  enforce_merge_condition          = false
}
```

## Argument Reference

* `project` - Required. Project Key that contains target repository.
* `repository` - Required. Repository slug of target repository.
* `allow_push_after_pr` - Optional. Enable users to push to branches that have an open outgoing pull requests. Default `true`.
* `unapprove_pr_after_source_change` - Optional. When code is pushed to a branch with an outgoing Pull Request, withdraw all approvals for this Pull Request.
* `unapprove_pr_after_target_change` - Optional. When code is pushed or merged to the target branch of a Pull Request, withdraw all approvals for this Pull Request. Default `false`.
* `enforce_merge_condition` - Optional. Enforce the Workzone merge condition for the project. Default `true`.
