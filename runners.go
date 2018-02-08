//
// Copyright 2017, Sander van Harmelen
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

package gitlab

import (
	"fmt"
	"net/url"
	"time"
)

// RunnersService handles communication with the runner related methods of the
// GitLab API.
//
// GitLab API docs: https://docs.gitlab.com/ce/api/runners.html
type RunnersService struct {
	client *Client
}

// Runner represents a GitLab CI Runner.
//
// GitLab API docs: https://docs.gitlab.com/ce/api/runners.html
type Runner struct {
	ID          int    `json:"id"`
	Description string `json:"description"`
	Active      bool   `json:"active"`
	IsShared    bool   `json:"is_shared"`
	Name        string `json:"name"`
	Online      bool   `json:"online"`
	Status      string `json:"status"`
}

// RunnersDetails represents a GitLab CI RunnerDetails.
//
// GitLab API docs: https://docs.gitlab.com/ce/api/runners.html
type RunnersDetails struct {
	Active       bool       `json:"active"`
	Architecture string     `json:"architecture"`
	Description  string     `json:"description"`
	ID           int        `json:"id"`
	IsShared     bool       `json:"is_shared"`
	ContactedAt  *time.Time `json:"contacted_at,omitempty"`
	Name         string     `json:"name"`
	Online       bool       `json:"online"`
	Status       string     `json:"status"`
	Platform     string     `json:"platform,omitempty"`
	Projects     []struct {
		ID                int    `json:"id"`
		Name              string `json:"name"`
		NameWithNamespace string `json:"name_with_namespace"`
		Path              string `json:"path"`
		PathWithNamespace string `json:"path_with_namespace"`
	} `json:"projects"`
	Token       string   `json:"Token"`
	Revision    string   `json:"revision,omitempty"`
	TagList     []string `json:"tag_list"`
	Version     string   `json:"version,omitempty"`
	AccessLevel string   `json:"access_level"`
}

// ListRunnersOptions represents the available ListRunners() options.
//
// GitLab API docs:
// https://docs.gitlab.com/ce/api/runners.html#list-owned-runners
type ListRunnersOptions struct {
	ListOptions
	Scope *string `url:"scope,omitempty" json:"scope,omitempty"`
}

// ListRunners gets a list of runners accessible by the authenticated user.
//
// GitLab API docs:
// https://docs.gitlab.com/ce/api/runners.html#list-owned-runners
func (s *RunnersService) ListRunners(opt *ListRunnersOptions, options ...OptionFunc) ([]*Runner, *Response, error) {
	req, err := s.client.NewRequest("GET", "runners", opt, options)
	if err != nil {
		return nil, nil, err
	}

	var rs []*Runner
	resp, err := s.client.Do(req, &rs)
	if err != nil {
		return nil, resp, err
	}

	return rs, resp, err
}

// ListAllRunners gets a list of all runners in the GitLab instance. Access is
// restricted to users with admin privileges.
//
// GitLab API docs:
// https://docs.gitlab.com/ce/api/runners.html#list-all-runners
func (s *RunnersService) ListAllRunners(opt *ListRunnersOptions, options ...OptionFunc) ([]*Runner, *Response, error) {
	req, err := s.client.NewRequest("GET", "runners/all", opt, options)
	if err != nil {
		return nil, nil, err
	}

	var rs []*Runner
	resp, err := s.client.Do(req, &rs)
	if err != nil {
		return nil, resp, err
	}

	return rs, resp, err
}

// GetRunnerDetails returns details for given runner.
//
// GitLab API docs:
// https://docs.gitlab.com/ce/api/runners.html#get-runner-39-s-details
func (s *RunnersService) GetRunnerDetails(rid interface{}, options ...OptionFunc) (*RunnersDetails, *Response, error) {
	runner, err := parseID(rid)
	if err != nil {
		return nil, nil, err
	}
	u := fmt.Sprintf("runners/%s", runner)

	req, err := s.client.NewRequest("GET", u, nil, options)
	if err != nil {
		return nil, nil, err
	}

	var rs *RunnersDetails
	resp, err := s.client.Do(req, &rs)
	if err != nil {
		return nil, resp, err
	}

	return rs, resp, err
}

// UpdateRunnersDetailsOptions represents the available UpdateRunnersDetails() options.
//
// GitLab API docs:
// https://docs.gitlab.com/ce/api/runners.html#update-runner-39-s-details
type UpdateRunnersDetailsOptions struct {
	Description *string  `url:"description,omitempty" json:"description,omitempty"`
	Active      *bool    `url:"active,omitempty" json:"active,omitempty"`
	TagList     []string `url:"tag_list[],omitempty" json:"tag_list,omitempty"`
	RunUntagged *bool    `url:"run_untagged,omitempty" json:"run_untagged,omitempty"`
	Locked      *bool    `url:"locked,omitempty" json:"locked,omitempty"`
	AccessLevel *string  `url:"access_level,omitempty" json:"access_level,omitempty"`
}

// UpdateRunnersDetails updates runners details
//
// GitLab API docs:
// https://docs.gitlab.com/ce/api/runners.html#update-runner-39-s-details
func (s *RunnersService) UpdateRunnersDetails(rid interface{}, opt *UpdateRunnersDetailsOptions, options ...OptionFunc) (*RunnersDetails, *Response, error) {
	runner, err := parseID(rid)
	if err != nil {
		return nil, nil, err
	}
	u := fmt.Sprintf("runners/%s", runner)

	req, err := s.client.NewRequest("PUT", u, opt, options)
	if err != nil {
		return nil, nil, err
	}

	var rs *RunnersDetails
	resp, err := s.client.Do(req, &rs)
	if err != nil {
		return nil, resp, err
	}

	return rs, resp, err
}

// RemoveARunner removes a runner
//
// GitLab API docs:
// https://docs.gitlab.com/ce/api/runners.html#remove-a-runner
func (s *RunnersService) RemoveARunner(rid interface{}, options ...OptionFunc) (*Response, error) {
	runner, err := parseID(rid)
	if err != nil {
		return nil, err
	}
	u := fmt.Sprintf("runners/%s", runner)

	req, err := s.client.NewRequest("DELETE", u, nil, options)
	if err != nil {
		return nil, err
	}

	return s.client.Do(req, nil)
}

// ListRunnersJobsOptions represents the available ListRunnersJobs()
// options. (one of running, success, failed, canceled)
//
// GitLab API docs:
// https://docs.gitlab.com/ce/api/runners.html#list-runner-39-s-jobs
type ListRunnersJobsOptions struct {
	ListOptions
	Status *BuildState `url:"status,omitempty" json:"status,omitempty"`
}

// ListRunnerJobs gets a list of jobs that are being processed or were processed by specified Runner.
//
// GitLab API docs:
// https://docs.gitlab.com/ce/api/runners.html#list-runner-39-s-jobs
func (s *RunnersService) ListRunnerJobs(rid interface{}, opt *ListRunnersJobsOptions, options ...OptionFunc) ([]*Job, *Response, error) {
	runner, err := parseID(rid)
	if err != nil {
		return nil, nil, err
	}
	u := fmt.Sprintf("runners/%s/jobs", runner)

	req, err := s.client.NewRequest("GET", u, opt, options)
	if err != nil {
		return nil, nil, err
	}

	var rs []*Job
	resp, err := s.client.Do(req, &rs)
	if err != nil {
		return nil, resp, err
	}

	return rs, resp, err
}

// ListProjectRunnersOptions represents the available ListProjectRunners()
// options.
//
// GitLab API docs:
// https://docs.gitlab.com/ce/api/runners.html#list-project-s-runners
type ListProjectRunnersOptions ListRunnersOptions

// ListProjectRunners gets a list of runners accessible by the authenticated user.
//
// GitLab API docs:
// https://docs.gitlab.com/ce/api/runners.html#list-project-s-runners
func (s *RunnersService) ListProjectRunners(pid interface{}, opt *ListProjectRunnersOptions, options ...OptionFunc) ([]*Runner, *Response, error) {
	project, err := parseID(pid)
	if err != nil {
		return nil, nil, err
	}
	u := fmt.Sprintf("projects/%s/runners", url.QueryEscape(project))

	req, err := s.client.NewRequest("GET", u, opt, options)
	if err != nil {
		return nil, nil, err
	}

	var rs []*Runner
	resp, err := s.client.Do(req, &rs)
	if err != nil {
		return nil, resp, err
	}

	return rs, resp, err
}

// EnableProjectRunnerOptions represents the available EnableProjectRunner()
// options.
//
// GitLab API docs:
// https://docs.gitlab.com/ce/api/runners.html#enable-a-runner-in-project
type EnableProjectRunnerOptions struct {
	RunnerID int `json:"runner_id"`
}

// EnableProjectRunner enables an available specific runner in the project.
//
// GitLab API docs:
// https://docs.gitlab.com/ce/api/runners.html#enable-a-runner-in-project
func (s *RunnersService) EnableProjectRunner(pid interface{}, opt *EnableProjectRunnerOptions, options ...OptionFunc) (*Runner, *Response, error) {
	project, err := parseID(pid)
	if err != nil {
		return nil, nil, err
	}
	u := fmt.Sprintf("projects/%s/runners", url.QueryEscape(project))

	req, err := s.client.NewRequest("POST", u, opt, options)
	if err != nil {
		return nil, nil, err
	}

	var r *Runner
	resp, err := s.client.Do(req, &r)
	if err != nil {
		return nil, resp, err
	}

	return r, resp, err
}

// DisableProjectRunner disables a specific runner from project.
//
// GitLab API docs:
// https://docs.gitlab.com/ce/api/runners.html#disable-a-runner-from-project
func (s *RunnersService) DisableProjectRunner(pid interface{}, rid interface{}, options ...OptionFunc) (*Response, error) {
	project, err := parseID(pid)
	if err != nil {
		return nil, err
	}

	runner, err := parseID(rid)
	if err != nil {
		return nil, err
	}

	u := fmt.Sprintf("projects/%s/runners/%s", url.QueryEscape(project), url.QueryEscape(runner))

	req, err := s.client.NewRequest("DELETE", u, nil, options)
	if err != nil {
		return nil, err
	}

	return s.client.Do(req, nil)
}