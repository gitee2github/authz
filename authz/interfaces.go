// Copyright (c) Huawei Technologies Co., Ltd. 2019. All rights reserved.
// authz is licensed under the Mulan PSL v2.
// You can use this software according to the terms and conditions of the Mulan PSL v2.
// You may obtain a copy of Mulan PSL v2 at:
//    http://license.coscl.org.cn/MulanPSL2
// THIS SOFTWARE IS PROVIDED ON AN "AS IS" BASIS, WITHOUT WARRANTIES OF ANY KIND, EITHER EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO NON-INFRINGEMENT, MERCHANTABILITY OR FIT FOR A PARTICULAR
// PURPOSE.
// See the Mulan PSL v2 for more details.
// Description: authorizer and auditor interface
// Author: zhangsong
// Create: 2019-01-18

package authz

import "github.com/docker/docker/pkg/authorization"

// Authorizer handles the authorization of isulad requests and responses
type Authorizer interface {
	Init() error
	LoadPolicies() error
	GetPolicies() []Policy
	AuthZRequest(req *authorization.Request) *authorization.Response
	AuthZResponse(req *authorization.Request) *authorization.Response
}

// Auditor audits the request and response sent from/to isulad daemon
type Auditor interface {
	AuditRequest(req *authorization.Request, resp *authorization.Response) error
	AuditResponse(req *authorization.Request, resp *authorization.Response) error
}
