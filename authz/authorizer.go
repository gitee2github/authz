// Copyright (c) Huawei Technologies Co., Ltd. 2019. All rights reserved.
// authz is licensed under the Mulan PSL v1.
// You can use this software according to the terms and conditions of the Mulan PSL v1.
// You may obtain a copy of Mulan PSL v1 at:
//    http://license.coscl.org.cn/MulanPSL
// THIS SOFTWARE IS PROVIDED ON AN "AS IS" BASIS, WITHOUT WARRANTIES OF ANY KIND, EITHER EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO NON-INFRINGEMENT, MERCHANTABILITY OR FIT FOR A PARTICULAR
// PURPOSE.
// See the Mulan PSL v1 for more details.
// Description: handle authz request and response for isulad
// Author: zhangsong
// Create: 2019-01-18

package authz

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/signal"
	"path"
	"regexp"
	"strings"
	"syscall"

	"github.com/docker/docker/pkg/authorization"
	"github.com/sirupsen/logrus"
)

// Policy is rbac policy
type Policy struct {
	Actions  []string `json:"actions"`  // Actions are the isulad actions
	Users    []string `json:"users"`    // Users are the users for which this policy apply to
	Name     string   `json:"name"`     // Name is the policy name
	Readonly bool     `json:"readonly"` // Readonly indicates this policy only allow get commands
}

type authorizer struct {
	policyPath string
	policies   []Policy
}

// NewAuthorizer creates a new authorizer
func NewAuthorizer(policyPath string) Authorizer {
	return &authorizer{policyPath: policyPath}
}

// Init loads the authz plugin configuration
func (f *authorizer) Init() error {
	err := f.LoadPolicies()
	if err != nil {
		return err
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGHUP)
	go func() {
		for range c {
			if err := f.LoadPolicies(); err != nil {
				logrus.Errorf("Error reloading policy %q", err.Error())
			}
		}
	}()

	return nil
}

func (f *authorizer) LoadPolicies() error {
	policies, err := parsePolicy(f.policyPath)
	if err != nil {
		return err
	}

	// Notify duplicate policies
	policyMap := make(map[string]string)
	for _, policy := range policies {
		for _, u := range policy.Users {
			if v, ok := policyMap[u]; ok {
				logrus.Warnf(
					"[policy: %s] User %q already appears in policy %q. Only single policy applies.",
					policy.Name,
					u,
					v,
				)
			}
			policyMap[u] = policy.Name
		}
	}
	f.policies = policies
	return nil
}

func parsePolicy(policyPath string) ([]Policy, error) {
	data, err := ioutil.ReadFile(path.Join(policyPath))
	if err != nil {
		return nil, err
	}

	var policies []Policy
	for _, line := range strings.Split(string(data), "\n") {
		if line == "" {
			continue
		}
		var policy Policy
		err := json.Unmarshal([]byte(line), &policy)
		if err != nil {
			logrus.Errorf("Failed to unmarshal policy %q %q", line, err.Error())
		}
		policies = append(policies, policy)
	}
	return policies, nil
}

func (f *authorizer) GetPolicies() []Policy {
	return f.policies
}

func (f *authorizer) AuthZRequest(request *authorization.Request) *authorization.Response {

	logrus.Debugf("Received AuthZ request, method: '%s', url: '%s'", request.RequestMethod, request.RequestURI)
	action := ParseRoute(request.RequestMethod, request.RequestURI)

	response := &authorization.Response{}

	for _, policy := range f.policies {
		for _, user := range policy.Users {
			if user == "" || user == request.User {
				for _, policyActionPattern := range policy.Actions {
					match, err := regexp.MatchString(policyActionPattern, action)
					if err != nil {
						logrus.Errorf(
							"Failed to recognize action %q against policy %q error %q",
							action,
							policyActionPattern,
							err.Error(),
						)
					}

					if match {
						if policy.Readonly && request.RequestMethod != "GET" {
							response.Allow = false
							response.Msg = fmt.Sprintf(
								"action '%s' not allowed for user '%s' by readonly policy %s",
								action,
								request.User,
								policy.Name,
							)
							return response
						}
						response.Allow = true
						return response
					}
				}
				response.Allow = false
				response.Msg = fmt.Sprintf(
					"action '%s' denied for user '%s' by policy '%s'",
					action,
					request.User,
					policy.Name,
				)
				return response
			}
		}
	}
	response.Allow = false
	response.Msg = fmt.Sprintf(
		"no policy applied (user: '%s' action: '%s')",
		request.User,
		action,
	)
	return response
}

func (f *authorizer) AuthZResponse(request *authorization.Request) *authorization.Response {
	return &authorization.Response{Allow: true}
}
