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
// https://bitbucket.org/izymessupport/workzone-public-apis/src/8fb698f4977575efe9b010dcd2ad105371a1a74a/workzone-api-generated.json#lines-1574
type WorkzoneReviewers struct {
	Project               string             `json:"projectKey,omitempty"`
	Repository            string             `json:"repoSlug,omitempty"`
	RefName               string             `json:"refName,omitempty"`
	Users                 []User             `json:"users,omitempty"`
	Groups                []string           `json:"groups,omitempty"`
	MandatoryUsers        []User             `json:"mandatoryUsers,omitempty"`
	MandatoryGroups       []string           `json:"mandatoryGroups,omitempty"`
	TopSuggestedReviewers int                `json:"topSuggestedReviewers,omitempty"`
	DaysInPast            int                `json:"daysInPast,omitempty"`
	FilePathReviewers     []FilePathReviewer `json:"filePathReviewers,omitempty"`
}

type FilePathReviewer struct {
	FilePathPattern       string   `json:"filePathPattern,omitempty"`
	IncludeExcludeEnabled bool     `json:"includeExcludeEnabled,omitempty"`
	ExcludeFilePaths      []string `json:"excludeFilePathPattern,omitempty"`
	Users                 []User   `json:"users,omitempty"`
	Groups                []string `json:"groups,omitempty"`
	MandatoryUsers        []User   `json:"mandatoryUsers,omitempty"`
	MandatoryGroups       []string `json:"mandatoryGroups,omitempty"`
}

func GetUserFromApiByUsername(userName string, client *BitbucketClient) (User, error) {
	req, err := client.Get(fmt.Sprintf("/rest/api/1.0/users/%s",
		url.PathEscape(userName),
	))
	var user User
	if req.StatusCode == 200 {
		body, _ := ioutil.ReadAll(req.Body)
		json.Unmarshal(body, &user)
		return user, nil
	} else {
		return user, fmt.Errorf("failed to find user %s: %+v", userName, err)
	}
}

// NEXT TWO METHODS ARE ALMOST THE SAME... CAN I DO SOMETHING WITH IT?...
// #1
func AddRepositoryReviewers(
	wz *WorkzoneReviewers,
	users []interface{},
	groups []interface{},
	mandatoryUsers []interface{},
	mandatoryGroups []interface{},
	client *BitbucketClient) error {
	// Add user reviewers to *WorkzoneReviewers
	for _, item := range users {
		user, err := GetUserFromApiByUsername(item.(string), client)
		if err != nil {
			return err
		}
		wz.Users = append(wz.Users, user)
	}
	// Add groups reviewers to *WorkzoneReviewers
	for _, item := range groups {
		wz.Groups = append(wz.Groups, item.(string))
	}

	// Add mandatory users reviewers to *WorkzoneReviewers
	for _, item := range mandatoryUsers {
		mandatoryUser, err := GetUserFromApiByUsername(item.(string), client)
		if err != nil {
			return err
		}
		wz.MandatoryUsers = append(wz.MandatoryUsers, mandatoryUser)
	}
	// Add mandatory groups reviewers to *WorkzoneReviewers
	for _, item := range mandatoryGroups {
		wz.MandatoryGroups = append(wz.MandatoryGroups, item.(string))
	}
	// wz.MandatoryGroups = stringArrayFromSchemaSet(mandatoryGroups)
	return nil
}

// #2
func AddFilepathReviewers(
	fpReviewers *FilePathReviewer,
	users []interface{},
	groups []interface{},
	mandatoryUsers []interface{},
	mandatoryGroups []interface{},
	client *BitbucketClient) error {
	// Add user reviewers to *FilePathReviewer
	for _, item := range users {
		user, err := GetUserFromApiByUsername(item.(string), client)
		if err != nil {
			return err
		}
		fpReviewers.Users = append(fpReviewers.Users, user)
	}
	// Add groups reviewers to *FilePathReviewer
	for _, item := range groups {
		fpReviewers.Groups = append(fpReviewers.Groups, item.(string))
	}
	// fpReviewers.Groups = stringArrayFromSchemaSet(groups)

	// Add mandatory users reviewers to *FilePathReviewer
	for _, item := range mandatoryUsers {
		mandatoryUser, err := GetUserFromApiByUsername(item.(string), client)
		if err != nil {
			return err
		}
		fpReviewers.MandatoryUsers = append(fpReviewers.MandatoryUsers, mandatoryUser)
	}
	// Add mandatory groups reviewers to *FilePathReviewer
	for _, item := range mandatoryGroups {
		fpReviewers.MandatoryGroups = append(fpReviewers.MandatoryGroups, item.(string))
	}
	// fpReviewers.MandatoryGroups = stringArrayFromSchemaSet(mandatoryGroups)
	return nil
}

func resourceWorkzoneReviewers() *schema.Resource {
	return &schema.Resource{
		Create: resourceWorkzoneReviewersCreate,
		Read:   resourceWorkzoneReviewersRead,
		Update: resourceWorkzoneReviewersCreate, // same as Create
		Delete: resourceWorkzoneReviewersDelete,
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
			"repository_reviewers_users": {
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Optional: true,
			},
			"repository_reviewers_groups": {
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Optional: true,
			},
			"mandatory_repository_reviewers_users": {
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Optional: true,
			},
			"mandatory_repository_reviewers_groups": {
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Optional: true,
			},
			"suggested_reviewers": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  0,
			},
			"suggested_reviewers_timespan": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  90,
			},
			"filepath_reviewers": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"filepath_pattern": {
							Type:     schema.TypeString,
							Required: true,
						},
						"include_exclude_enabled": {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  false,
						},
						"exclude_filepaths": {
							Type:     schema.TypeList,
							Elem:     &schema.Schema{Type: schema.TypeString},
							Optional: true,
						},
						"filepath_reviewers_users": {
							Type:     schema.TypeList,
							Elem:     &schema.Schema{Type: schema.TypeString},
							Optional: true,
						},
						"filepath_reviewers_groups": {
							Type:     schema.TypeList,
							Elem:     &schema.Schema{Type: schema.TypeString},
							Optional: true,
						},
						"mandatory_filepath_reviewers_users": {
							Type:     schema.TypeList,
							Elem:     &schema.Schema{Type: schema.TypeString},
							Optional: true,
						},
						"mandatory_filepath_reviewers_groups": {
							Type:     schema.TypeList,
							Elem:     &schema.Schema{Type: schema.TypeString},
							Optional: true,
						},
					},
				},
			},
		},
	}
}

func newWorkzoneReviewersFromResource(d *schema.ResourceData, client *BitbucketClient) *WorkzoneReviewers {
	reviewers := &WorkzoneReviewers{
		Project:               d.Get("project").(string),
		Repository:            d.Get("repository").(string),
		RefName:               d.Get("refname").(string),
		TopSuggestedReviewers: d.Get("suggested_reviewers").(int),
		DaysInPast:            d.Get("suggested_reviewers_timespan").(int),
	}

	AddRepositoryReviewers(
		reviewers,
		d.Get("repository_reviewers_users").([]interface{}),
		d.Get("repository_reviewers_groups").([]interface{}),
		d.Get("mandatory_repository_reviewers_users").([]interface{}),
		d.Get("mandatory_repository_reviewers_groups").([]interface{}),
		client,
	)

	for _, item := range d.Get("filepath_reviewers").([]interface{}) {
		fpReviewer := &FilePathReviewer{
			FilePathPattern:       item.(map[string]interface{})["filepath_pattern"].(string),
			IncludeExcludeEnabled: item.(map[string]interface{})["include_exclude_enabled"].(bool),
		}

		for _, item := range item.(map[string]interface{})["exclude_filepaths"].([]interface{}) {
			fpReviewer.ExcludeFilePaths = append(fpReviewer.ExcludeFilePaths, item.(string))
		}

		AddFilepathReviewers(
			fpReviewer,
			item.(map[string]interface{})["filepath_reviewers_users"].([]interface{}),
			item.(map[string]interface{})["filepath_reviewers_groups"].([]interface{}),
			item.(map[string]interface{})["mandatory_filepath_reviewers_users"].([]interface{}),
			item.(map[string]interface{})["mandatory_filepath_reviewers_groups"].([]interface{}),
			client,
		)

		reviewers.FilePathReviewers = append(reviewers.FilePathReviewers, *fpReviewer)
	}

	return reviewers
}

func resourceWorkzoneReviewersCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*BitbucketServerProvider).BitbucketClient
	wz := newWorkzoneReviewersFromResource(d, client)

	bytedata, err := json.Marshal(wz)

	if err != nil {
		return err
	}

	// When you install the Workzone plugin, then create a repository and immediately after that try to update the Workzone settings,
	// the API can return 404. That's why the POST call is wrapped into the retry function.
	err = resource.Retry(time.Minute,
		func() *resource.RetryError {
			_, err = client.Post(fmt.Sprintf("/rest/workzoneresource/1.0/branch/reviewers/%s/%s",
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

	return resourceWorkzoneReviewersRead(d, m)
}

func resourceWorkzoneReviewersRead(d *schema.ResourceData, m interface{}) error {
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
	req, err := client.Get(fmt.Sprintf("/rest/workzoneresource/1.0/branch/reviewers/%s/%s",
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

	var wz []WorkzoneReviewers

	body, readerr := ioutil.ReadAll(req.Body)
	if readerr != nil {
		return readerr
	}

	decodeerr := json.Unmarshal(body, &wz)
	if decodeerr != nil {
		return decodeerr
	}

	if len(wz) > 0 {
		// Workzone API returns array of items [{...}]
		d.Set("repository_reviewers_users", wz[0].Users)
		d.Set("repository_reviewers_groups", wz[0].Groups)
		d.Set("mandatory_repository_reviewers_users", wz[0].MandatoryUsers)
		d.Set("mandatory_repository_reviewers_groups", wz[0].MandatoryGroups)
		d.Set("suggested_reviewers", wz[0].TopSuggestedReviewers)
		d.Set("suggested_reviewers_timespan", wz[0].DaysInPast)
		d.Set("filepath_reviewers", wz[0].FilePathReviewers)
	} else {
		log.Printf("[WARN] Workzone Reviewers object (%s) is empty, removing from state", d.Id())
		d.SetId("")
		return nil
	}

	return nil
}

func resourceWorkzoneReviewersDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*BitbucketServerProvider).BitbucketClient

	project := d.Get("project").(string)
	repository := d.Get("repository").(string)

	// For some reason Workzone API requires DELETE call to contain body
	wz := newWorkzoneReviewersFromResource(d, client)
	bytedata, err := json.Marshal(wz)

	if err != nil {
		return err
	}

	_, err = client.DeleteWithBody(fmt.Sprintf("/rest/workzoneresource/1.0/branch/reviewers/%s/%s",
		url.QueryEscape(project),
		url.QueryEscape(repository),
	), bytes.NewBuffer(bytedata))

	return err
}
