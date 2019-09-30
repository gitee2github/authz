// Copyright (c) Huawei Technologies Co., Ltd. 2019. All rights reserved.
// authz is licensed under the Mulan PSL v1.
// You can use this software according to the terms and conditions of the Mulan PSL v1.
// You may obtain a copy of Mulan PSL v1 at:
//    http://license.coscl.org.cn/MulanPSL
// THIS SOFTWARE IS PROVIDED ON AN "AS IS" BASIS, WITHOUT WARRANTIES OF ANY KIND, EITHER EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO NON-INFRINGEMENT, MERCHANTABILITY OR FIT FOR A PARTICULAR
// PURPOSE.
// See the Mulan PSL v1 for more details.
// Description: handle isulad http request
// Author: zhangsong
// Create: 2019-01-18

package core

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"

	"github.com/docker/docker/pkg/authorization"
	"github.com/docker/docker/pkg/plugins"
	"github.com/sirupsen/logrus"
)

// HandleFunc handle function for authz
type HandleFunc func(w http.ResponseWriter, r *http.Request)

// HandleActive handle authz plugin active
func (a *AuthZServer) HandleActive() HandleFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		b, err := json.Marshal(plugins.Manifest{Implements: []string{authorization.AuthZApiImplements}})
		if err != nil {
			writeErr(w, err)
			return
		}
		_, err = w.Write(b)
		if err != nil {
			logrus.Warnf("write http response err:%v", err)
		}
	}
}

// HandleRequest handle authz request
func (a *AuthZServer) HandleRequest() HandleFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			writeErr(w, err)
			return
		}

		req := &authorization.Request{}
		err = json.Unmarshal(body, req)
		if err != nil {
			writeErr(w, err)
			return
		}

		resp := a.authorizer.AuthZRequest(req)
		if resp != nil {
			logrus.Debugf(resp.Msg)
		}

		err = a.auditor.AuditRequest(req, resp)
		if err != nil {
			logrus.Errorf("Failed to audit request '%v'", err)
		}

		writeResponse(w, resp)
	}
}

// HandleResponse handle authz response
func (a *AuthZServer) HandleResponse() HandleFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()
		body, err := ioutil.ReadAll(r.Body)

		if err != nil {
			writeErr(w, err)
			return
		}

		req := &authorization.Request{}
		err = json.Unmarshal(body, req)

		if err != nil {
			writeErr(w, err)
			return
		}

		resp := a.authorizer.AuthZResponse(req)
		err = a.auditor.AuditResponse(req, resp)
		if err != nil {
			logrus.Errorf("Failed to audit response '%v'", err)
		}

		writeResponse(w, resp)
	}
}

func writeErr(w http.ResponseWriter, err error) {
	writeResponse(w, &authorization.Response{Err: err.Error()})
}

func writeResponse(w http.ResponseWriter, resp *authorization.Response) {
	data, err := json.Marshal(resp)
	if err != nil {
		logrus.Errorf("Failed to marshal authz response %q", err.Error())
	} else {
		_, err = w.Write(data)
		if err != nil {
			logrus.Warnf("write http response err:%v", err)
		}
	}

	if resp == nil || resp.Err != "" {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

// HandleIsuladRequest handle isulad authz request
func (a *AuthZServer) HandleIsuladRequest() HandleFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if r.Body != nil {
				r.Body.Close()
			}
		}()
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			logrus.Errorf("Failed to read body from request: %s", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		items := strings.Split(string(body), ":")
		if len(items) != 2 { // Standard format length
			logrus.Errorf("Bad format: %v", items)
			w.WriteHeader(http.StatusBadRequest)
		} else {
			username := items[0]
			action := items[1]
			respCode := a.AuthIsuladUser(username, action)
			w.WriteHeader(respCode)
		}
	}
}

// AuthIsuladUser authorize user for isulad http request
func (a *AuthZServer) AuthIsuladUser(username string, action string) int {
	if err := a.authorizer.LoadPolicies(); err != nil {
		logrus.Errorf("Failed to load policies: %s", err)
		return http.StatusInternalServerError
	}
	for _, policy := range a.authorizer.GetPolicies() {
		for _, user := range policy.Users {
			if user == "" || user == username {
				for _, policyActionPattern := range policy.Actions {
					match, err := regexp.MatchString(policyActionPattern, action)
					if err != nil {
						logrus.Errorf(
							"Failed to recognize action %q against policy %q error %q",
							action,
							policyActionPattern,
							err.Error(),
						)
						return http.StatusInternalServerError
					}

					if match {
						if policy.Readonly {
							logrus.Errorf(
								"action '%s' not allowed for user '%s' by readonly policy %s",
								action,
								username,
								policy.Name,
							)
							return http.StatusForbidden
						}
						return http.StatusOK
					}
				}
				logrus.Errorf(
					"action '%s' denied for user '%s' by policy '%s'",
					action,
					username,
					policy.Name,
				)
				return http.StatusForbidden
			}
		}
	}
	logrus.Errorf(
		"no policy applied (user: '%s' action: '%s')",
		username,
		action,
	)
	return http.StatusNotFound
}
