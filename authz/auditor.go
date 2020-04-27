// Copyright (c) Huawei Technologies Co., Ltd. 2019. All rights reserved.
// authz is licensed under the Mulan PSL v2.
// You can use this software according to the terms and conditions of the Mulan PSL v2.
// You may obtain a copy of Mulan PSL v2 at:
//    http://license.coscl.org.cn/MulanPSL2
// THIS SOFTWARE IS PROVIDED ON AN "AS IS" BASIS, WITHOUT WARRANTIES OF ANY KIND, EITHER EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO NON-INFRINGEMENT, MERCHANTABILITY OR FIT FOR A PARTICULAR
// PURPOSE.
// See the Mulan PSL v2 for more details.
// Description: audit request or response directly to standard output,
//              available for isulad
// Author: zhangsong
// Create: 2019-01-18

package authz

import (
	"fmt"
	"log/syslog"

	"github.com/docker/docker/pkg/authorization"
	"github.com/sirupsen/logrus"
	logrus_syslog "github.com/sirupsen/logrus/hooks/syslog"
)

// auditor audit request/response directly to standard output
type auditor struct {
	logger *logrus.Logger
}

// NewAuditor returns a new authz auditor
func NewAuditor() Auditor {
	return &auditor{}
}

func (b *auditor) AuditRequest(req *authorization.Request, resp *authorization.Response) error {

	if req == nil || resp == nil {
		return fmt.Errorf("Authorization request or response is nil")
	}

	err := b.init()
	if err != nil {
		return err
	}

	fields := logrus.Fields{
		"method": req.RequestMethod,
		"uri":    req.RequestURI,
		"user":   req.User,
		"allow":  resp.Allow,
		"msg":    resp.Msg,
	}

	if resp != nil || resp.Err != "" {
		fields["err"] = resp.Err
	}

	b.logger.WithFields(fields).Info("Request")
	return nil
}

func (b *auditor) AuditResponse(req *authorization.Request, resp *authorization.Response) error {
	return nil
}

// init inits the auditor logger
func (b *auditor) init() error {
	if b.logger != nil {
		return nil
	}

	b.logger = logrus.New()
	b.logger.Formatter = &logrus.JSONFormatter{}
	hook, err := logrus_syslog.NewSyslogHook("", "", syslog.LOG_ERR, "authz")
	if err != nil {
		return err
	}
	b.logger.Hooks.Add(hook)
	return nil
}
