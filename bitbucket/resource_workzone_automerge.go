package bitbucket

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	"io/ioutil"

	"net/http"
	"net/url"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
)

// The struct represents this JSON payload:
// https://bitbucket.org/izymessupport/workzone-public-apis/src/8fb698f4977575efe9b010dcd2ad105371a1a74a/workzone-api-generated.json#lines-1688
type WorkzoneAutoMerge struct {
	Project                             string `json:"projectKey,omitempty"`
	Repository                          string `json:"repoSlug,omitempty"`
	RefName                             string `json:"refName,omitempty"`
	SrcRefName                          string `json:"srcRefName,omitempty"`
	MergeCondition                      string `json:"mergeCondition,omitempty"`
	MergeStrategyId                     string `json:"mergeStrategyId,omitempty"`
	ApprovalQuota                       string `json:"approvalQuota,omitempty"`
	ApprovalCount                       int    `json:"approvalCount,omitempty"`
	MandatoryApprovalCount              int    `json:"mandatoryApprovalCount,omitempty"`
	GroupQuota                          int    `json:"groupQuota,omitempty"`
	AutoMergeEnabled                    bool   `json:"autoMergeEnabled,omitempty"`
	ApprovalQuotaEnabled                bool   `json:"approvalQuotaEnabled,omitempty"`
	DeleteSourceBranch                  bool   `json:"deleteSourceBranch,omitempty"`
	WatchBuildResult                    bool   `json:"watchBuildResult,omitempty"`
	WatchTaskCompletion                 bool   `json:"watchTaskCompletion,omitempty"`
	RequiredSignatures                  int    `json:"requiredSignaturesCount,omitempty"`
	RequiredBuilds                      int    `json:"requiredBuildsCount,omitempty"`
	IgnoreContributingReviewersApproval bool   `json:"ignoreContributingReviewersApproval,omitempty"`
	EnableNeedsWorkVeto                 bool   `json:"enableNeedsWorkVeto,omitempty"`
	AutomergeUsers                      []User `json:"automergeUsers,omitempty"`
}

func resourceWorkzoneAutoMerge() *schema.Resource {
	return &schema.Resource{
		Create: resourceWorkzoneAutoMergeCreate,
		Read:   resourceWorkzoneAutoMergeRead,
		Update: resourceWorkzoneAutoMergeCreate, // same as Create
		Delete: resourceWorkzoneAutoMergeDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"project": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"repository": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"refname": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"source_refname": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"merge_condition": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"merge_strategy_id": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "none-inherit",
			},
			"approval_quota_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"approval_quota": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"approval_count": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"mandatory_approval_count": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"group_quota": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"delete_source_branch": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"watch_build_result": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"watch_task_completion": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"required_signatures": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"required_builds": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"ignore_contributing_reviewers_approval": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"enable_needs_work_veto": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"automerge_users": {
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Optional: true,
			},
		},
	}
}

func newWorkzoneAutoMergeFromResource(d *schema.ResourceData, client *BitbucketClient) *WorkzoneAutoMerge {
	autoMergeConfig := &WorkzoneAutoMerge{
		Project:                             d.Get("project").(string),
		Repository:                          d.Get("repository").(string),
		RefName:                             d.Get("refname").(string),
		SrcRefName:                          d.Get("source_refname").(string),
		MergeCondition:                      d.Get("merge_condition").(string),
		MergeStrategyId:                     d.Get("merge_strategy_id").(string),
		ApprovalQuota:                       d.Get("approval_quota").(string),
		ApprovalCount:                       d.Get("approval_count").(int),
		MandatoryApprovalCount:              d.Get("mandatory_approval_count").(int),
		GroupQuota:                          d.Get("group_quota").(int),
		ApprovalQuotaEnabled:                d.Get("approval_quota_enabled").(bool),
		DeleteSourceBranch:                  d.Get("delete_source_branch").(bool),
		WatchBuildResult:                    d.Get("watch_build_result").(bool),
		WatchTaskCompletion:                 d.Get("watch_task_completion").(bool),
		RequiredSignatures:                  d.Get("required_signatures").(int),
		RequiredBuilds:                      d.Get("required_builds").(int),
		IgnoreContributingReviewersApproval: d.Get("ignore_contributing_reviewers_approval").(bool),
		EnableNeedsWorkVeto:                 d.Get("enable_needs_work_veto").(bool),
		AutomergeUsers:                      []User{},
	}

	for _, item := range d.Get("automerge_users").([]interface{}) {
		mandatoryUser, _ := GetUserFromApiByUsername(item.(string), client)
		autoMergeConfig.AutomergeUsers = append(autoMergeConfig.AutomergeUsers, mandatoryUser)
	}

	return autoMergeConfig
}

func resourceWorkzoneAutoMergeCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*BitbucketServerProvider).BitbucketClient
	wz := newWorkzoneAutoMergeFromResource(d, client)

	bytedata, err := json.Marshal(wz)

	if err != nil {
		return err
	}

	// When you install the Workzone plugin, then create a repository and immediately after that try to update the Workzone settings,
	// the API can return 404. That's why the POST call is wrapped into the retry function.
	err = resource.Retry(time.Minute,
		func() *resource.RetryError {
			_, err = client.Post(fmt.Sprintf("/rest/workzoneresource/1.0/branch/automerge/%s/%s",
				wz.Project,
				wz.Repository,
			), bytes.NewBuffer(bytedata))
			if err != nil {
				return resource.RetryableError(fmt.Errorf("waiting for workzone settings to become available"))
			} else {
				return nil
			}
		})
	if err != nil {
		return err
	}

	d.SetId(fmt.Sprintf("%s|%s|%s", wz.Project, wz.Repository, wz.RefName))

	return resourceWorkzoneAutoMergeRead(d, m)
}

func resourceWorkzoneAutoMergeRead(d *schema.ResourceData, m interface{}) error {
	id := d.Id()
	if id != "" {
		idparts := strings.Split(id, "|")
		if len(idparts) == 3 {
			_ = d.Set("project", idparts[0])
			_ = d.Set("repository", idparts[1])
			_ = d.Set("refname", idparts[2])
		} else {
			return fmt.Errorf("incorrect ID format, should match `project|repository|refname`")
		}
	}

	project := d.Get("project").(string)
	repository := d.Get("repository").(string)

	client := m.(*BitbucketServerProvider).BitbucketClient
	req, err := client.Get(fmt.Sprintf("/rest/workzoneresource/1.0/branch/automerge/%s/%s",
		url.PathEscape(project),
		url.PathEscape(repository),
	))

	if err != nil {
		return err
	}

	if req.StatusCode == http.StatusNotFound {
		log.Printf("[WARN] Workzone Reviewers object (%s) not found, removing from state", d.Id())
		d.SetId("")
		return nil
	}

	var wz []WorkzoneAutoMerge

	body, readerr := ioutil.ReadAll(req.Body)
	if readerr != nil {
		return readerr
	}

	decodeerr := json.Unmarshal(body, &wz)
	if decodeerr != nil {
		return decodeerr
	}

	// Workzone API returns array of items [{...}]
	if len(wz) > 0 {
		d.Set("source_refname", wz[0].SrcRefName)
		d.Set("merge_condition", wz[0].MergeCondition)
		d.Set("merge_strategy_id", wz[0].MergeStrategyId)
		d.Set("approval_quota_enabled", wz[0].ApprovalQuotaEnabled)
		d.Set("approval_quota", wz[0].ApprovalQuota)
		d.Set("approval_count", wz[0].ApprovalCount)
		d.Set("mandatory_approval_count", wz[0].MandatoryApprovalCount)
		d.Set("group_quota", wz[0].GroupQuota)
		d.Set("delete_source_branch", wz[0].DeleteSourceBranch)
		d.Set("watch_build_result", wz[0].WatchBuildResult)
		d.Set("required_signatures", wz[0].RequiredSignatures)
		d.Set("required_builds", wz[0].RequiredBuilds)
		d.Set("ignore_contributing_reviewers_approval", wz[0].IgnoreContributingReviewersApproval)
		d.Set("enable_needs_work_veto", wz[0].EnableNeedsWorkVeto)
		d.Set("automerge_users", wz[0].AutomergeUsers)
	} else {
		log.Printf("[WARN] Workzone AutoMerge object (%s) is empty, removing from state", d.Id())
		d.SetId("")
		return nil
	}

	return nil
}

func resourceWorkzoneAutoMergeDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*BitbucketServerProvider).BitbucketClient

	project := d.Get("project").(string)
	repository := d.Get("repository").(string)

	// For some reason Workzone API requires DELETE call to contain body
	// Content of the body doesn't matter
	wz := newWorkzoneAutoMergeFromResource(d, client)
	bytedata, err := json.Marshal(wz)

	if err != nil {
		return err
	}

	_, err = client.DeleteWithBody(fmt.Sprintf("/rest/workzoneresource/1.0/branch/automerge/%s/%s",
		url.QueryEscape(project),
		url.QueryEscape(repository),
	), bytes.NewBuffer(bytedata))

	return err
}
