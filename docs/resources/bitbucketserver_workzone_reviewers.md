# Resource: bitbucketserver_workzone_reviewers

Provides the ability to manage Workzone reviewers settings.

## Example Usage

```hcl
resource "bitbucketserver_workzone_reviewers" "test" {
  project                               = "MYPROJ"
  repository                            = "repo"
  refname                               = "refs/heads/master"
  repository_reviewers_users            = ["admin"]
  repository_reviewers_groups           = ["group_1"]
  mandatory_repository_reviewers_users  = ["admin"]
  mandatory_repository_reviewers_groups = ["group_1"]
  suggested_reviewers                   = 2
  suggested_reviewers_timespan          = 100
  filepath_reviewers {
    filepath_pattern                    = "**/path_1"
    filepath_reviewers_users            = ["admin"]
    filepath_reviewers_groups           = ["group_1"]
    mandatory_filepath_reviewers_users  = ["admin"]
    mandatory_filepath_reviewers_groups = ["group_1"]
    include_exclude_enabled             = true
    exclude_filepaths                   = ["**/path_1/*.tfstate"]
  }
  filepath_reviewers {
    filepath_pattern                    = "**/path_2"
    filepath_reviewers_users            = ["admin"]
    filepath_reviewers_groups           = ["group_1"]
    mandatory_filepath_reviewers_users  = ["admin"]
    mandatory_filepath_reviewers_groups = ["group_1"]
    include_exclude_enabled             = false
    exclude_filepaths                   = ["**/path_2/*.tfstate", "**/path_2/main.tf"]
  }
}
```

## Argument Reference

* `project` - Required. Project Key that contains target repository.
* `repository` - Required. Repository slug of target repository.
* `refname` - Required. Refname of the destination branch for the PRs.
* `repository_reviewers_users` - Optional. List of users added as reviewers for the PRs to the target branch.
* `repository_reviewers_groups` - Optional. List of groups of users added as reviewers for the PRs to the target branch.
* `mandatory_repository_reviewers_users` - Optional. List of users added as mandatory reviewers for the PRs to the target branch. PR can't be merged unless **ALL** mandatory users approve it.
* `mandatory_repository_reviewers_groups` - Optional. List of groups of users added as mandatory reviewers for the PRs to the target branch. PR can't be merged unless **ALL** mandatory users approve it.
* `suggested_reviewers` - Optional. Controls the number of committers added as suggested reviewers to the PRs. Default `0`.
* `suggested_reviewers_timespan` - Optional. Controls the time interval from which the `suggested_reviewers` are populated. Default `90`.
* `filepath_reviewers` - Optional. Reviewers settings for specific file paths in the repository.
* `filepath_reviewers.filepath_pattern` - Required. Pattern for the file paths in the repository.
* `filepath_reviewers.include_exclude_enabled` - Optional. Set to `true` if you want to exclude some sub-paths from the `filepath_reviewers.filepath_pattern`. Default `false`.
* `filepath_reviewers.exclude_filepaths` - Optional. List of file paths' patterns to exclude.
* `filepath_reviewers.filepath_reviewers_users` - Optional. List of users added as reviewers for the context path.
* `filepath_reviewers.filepath_reviewers_groups` - Optional. List of groups of users added as reviewers for the context path.
* `filepath_reviewers.mandatory_filepath_reviewers_users` - Optional. List of users added as mandatory reviewers for the context path. PR can't be merged unless **ALL** mandatory users approve it.
* `filepath_reviewers.mandatory_filepath_reviewers_users` - Optional. List of groups of users added as mandatory reviewers for the context path. PR can't be merged unless **ALL** mandatory users approve it.
