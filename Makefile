# Copyright (c) Huawei Technologies Co., Ltd. 2019. All rights reserved.
# authz is licensed under the Mulan PSL v1.
# You can use this software according to the terms and conditions of the Mulan PSL v1.
# You may obtain a copy of Mulan PSL v1 at:
#    http://license.coscl.org.cn/MulanPSL
# THIS SOFTWARE IS PROVIDED ON AN "AS IS" BASIS, WITHOUT WARRANTIES OF ANY KIND, EITHER EXPRESS OR
# IMPLIED, INCLUDING BUT NOT LIMITED TO NON-INFRINGEMENT, MERCHANTABILITY OR FIT FOR A PARTICULAR
# PURPOSE.
# See the Mulan PSL v1 for more details.
# Description: authz makefile
# Author: zhangsong
# Create: 2019-01-18

VERSION ?= v1.0.0

default: binary

ENV = CGO_ENABLED=0
GO_LDFLAGS = "-X main.version=$(VERSION)"
GOMOD = "-mod=vendor"

binary:
	mkdir -p bin/
	$(ENV) go build $(GOMOD) -o bin/authz-broker --ldflags $(GO_LDFLAGS) -a -installsuffix cgo ./main.go

clean:
	rm -rf bin/
