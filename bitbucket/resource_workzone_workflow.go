package bitbucket

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
)

// The struct represents this JSON payload:
// https://bitbucket.org/izymessupport/workzone-public-apis/src/8fb698f4977575efe9b010dcd2ad105371a1a74a/workzone-api-generated.json#lines-1767
type WorkzoneWorkflow struct {
	Project                 string `json:"projectKey,omitempty"`
	Repository              string `json:"repoSlug,omitempty"`
	PushAfterPR             bool   `json:"pushAfterPullReq"`
	UnapprovePRSourceChange bool   `json:"unapprovePullReq,omitempty"`
	UnapprovePRTargetChange bool   `json:"unapprovePullReqTargetRefChange,omitempty"`
	EnforceMergeCondition   bool   `json:"enableMergeConditionVeto,omitempty"`
}

func resourceWorkzoneWorkflow() *schema.Resource {
	return &schema.Resource{
		Create: resourceWorkzoneWorkflowCreate,
		Read:   resourceWorkzoneWorkflowRead,
		Update: resourceWorkzoneWorkflowCreate, // same as Create
		Delete: resourceWorkzoneWorkflowDelete,
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
			"allow_push_after_pr": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"unapprove_pr_after_source_change": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"unapprove_pr_after_target_change": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"enforce_merge_condition": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
		},
	}
}

func resourceWorkzoneWorkflowCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*BitbucketServerProvider).BitbucketClient
	wz := newWorkzoneWorkflowFromResource(d, client)

	bytedata, err := json.Marshal(wz)

	if err != nil {
		return err
	}

	// When you install the Workzone plugin, then create a repository and immediately after that try to update the Workzone settings,
	// the API can return 404. That's why the POST call is wrapped into the retry function.
	err = resource.Retry(time.Minute,
		func() *resource.RetryError {
			_, err = client.Post(fmt.Sprintf("/rest/workzoneresource/1.0/workflow/%s/%s",
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

	d.SetId(fmt.Sprintf("%s|%s", wz.Project, wz.Repository))

	return resourceWorkzoneWorkflowRead(d, m)
}

func newWorkzoneWorkflowFromResource(d *schema.ResourceData, client *BitbucketClient) *WorkzoneWorkflow {
	workflow := &WorkzoneWorkflow{
		Project:                 d.Get("project").(string),
		Repository:              d.Get("repository").(string),
		PushAfterPR:             d.Get("allow_push_after_pr").(bool),
		UnapprovePRSourceChange: d.Get("unapprove_pr_after_source_change").(bool),
		UnapprovePRTargetChange: d.Get("unapprove_pr_after_target_change").(bool),
		EnforceMergeCondition:   d.Get("enforce_merge_condition").(bool),
	}

	return workflow
}

func resourceWorkzoneWorkflowRead(d *schema.ResourceData, m interface{}) error {
	id := d.Id()

	if id != "" {
		idparts := strings.Split(id, "|")
		if len(idparts) == 2 {
			_ = d.Set("project", idparts[0])
			_ = d.Set("repository", idparts[1])
		} else {
			return fmt.Errorf("incorrect ID format, should match `project|repository`")
		}
	}

	client := m.(*BitbucketServerProvider).BitbucketClient
	req, err := client.Get(fmt.Sprintf("/rest/workzoneresource/1.0/workflow/%s/%s",
		url.PathEscape(d.Get("project").(string)),
		url.PathEscape(d.Get("repository").(string)),
	))

	if err != nil {
		return err
	}

	if req.StatusCode == http.StatusNotFound {
		log.Printf("[WARN] Workzone Reviewers object (%s) not found, removing from state", d.Id())
		d.SetId("")
		return nil
	}

	var wz WorkzoneWorkflow

	body, readerr := ioutil.ReadAll(req.Body)
	if readerr != nil {
		return readerr
	}

	decodeerr := json.Unmarshal(body, &wz)
	if decodeerr != nil {
		return decodeerr
	}

	d.Set("allow_push_after_pr", wz.PushAfterPR)
	d.Set("unapprove_pr_after_source_change", wz.UnapprovePRSourceChange)
	d.Set("unapprove_pr_after_target_change", wz.UnapprovePRTargetChange)
	d.Set("enforce_merge_condition", wz.EnforceMergeCondition)

	return nil
}

func resourceWorkzoneWorkflowDelete(d *schema.ResourceData, m interface{}) error {
	project := d.Get("project").(string)
	repository := d.Get("repository").(string)
	client := m.(*BitbucketServerProvider).BitbucketClient
	_, err := client.Delete(fmt.Sprintf("/rest/workzoneresource/1.0/workflow/%s/%s",
		url.QueryEscape(project),
		url.QueryEscape(repository),
	))

	return err
}
