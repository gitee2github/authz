// Copyright (c) Huawei Technologies Co., Ltd. 2018-2019. All rights reserved.
// authz is licensed under the Mulan PSL v2.
// You can use this software according to the terms and conditions of the Mulan PSL v2.
// You may obtain a copy of Mulan PSL v2 at:
//    http://license.coscl.org.cn/MulanPSL2
// THIS SOFTWARE IS PROVIDED ON AN "AS IS" BASIS, WITHOUT WARRANTIES OF ANY KIND, EITHER EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO NON-INFRINGEMENT, MERCHANTABILITY OR FIT FOR A PARTICULAR
// PURPOSE.
// See the Mulan PSL v2 for more details.
// Description: authz server
// Author: liruilin
// Create: 2018-04-26

package core

import (
	"fmt"
	"net"
	"net/http"
	"os"

	"github.com/docker/docker/pkg/authorization"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"isula.org/authz/authz"
)

const (
	authZPluginName   = "authz-broker"
	authZPluginFolder = "/run/isulad/plugins"
)

// AuthZServer describes the  authz server
type AuthZServer struct {
	authorizer authz.Authorizer
	auditor    authz.Auditor
	listener   net.Listener
}

// NewAuthZServer creates a new authorization server
func NewAuthZServer(authorizer authz.Authorizer, auditor authz.Auditor) *AuthZServer {
	return &AuthZServer{authorizer: authorizer, auditor: auditor}
}

// Start starts the authorization server
func (a *AuthZServer) Start() error {

	err := a.authorizer.Init()

	if err != nil {
		return err
	}

	if _, err := os.Stat(authZPluginFolder); os.IsNotExist(err) {
		logrus.Infof("Creating authz plugin folder %q", authZPluginFolder)
		err = os.MkdirAll("/run/isulad/plugins/", 0750)
		if err != nil {
			return err
		}
	}

	pluginPath := fmt.Sprintf("%s/%s.sock", authZPluginFolder, authZPluginName)

	if err := os.Remove(pluginPath); err != nil {
		if !os.IsNotExist(err) {
			logrus.Errorf("Failed to remove pluginPath err: %s", err.Error())
		}
	}
	a.listener, err = net.ListenUnix("unix", &net.UnixAddr{Name: pluginPath, Net: "unix"})
	if err != nil {
		return err
	}
	if err := os.Chmod(pluginPath, 0660); err != nil {
		logrus.Errorf("Failed to chmod pluginPath err: %s", err.Error())
	}

	return a.route()
}

func (a *AuthZServer) route() error {
	router := mux.NewRouter()
	router.HandleFunc("/Plugin.Activate", a.HandleActive())
	router.HandleFunc(fmt.Sprintf("/%s", authorization.AuthZApiRequest), a.HandleRequest())
	router.HandleFunc(fmt.Sprintf("/%s", authorization.AuthZApiResponse), a.HandleResponse())
	router.HandleFunc("/isulad.auth", a.HandleIsuladRequest())
	return http.Serve(a.listener, router)
}

// Stop stop the authorization server
func (a *AuthZServer) Stop() error {
	if a.listener == nil {
		return fmt.Errorf("Listener is nil")
	}
	return a.listener.Close()
}
