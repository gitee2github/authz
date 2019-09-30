// Copyright (c) Huawei Technologies Co., Ltd. 2019. All rights reserved.
// authz is licensed under the Mulan PSL v1.
// You can use this software according to the terms and conditions of the Mulan PSL v1.
// You may obtain a copy of Mulan PSL v1 at:
//    http://license.coscl.org.cn/MulanPSL
// THIS SOFTWARE IS PROVIDED ON AN "AS IS" BASIS, WITHOUT WARRANTIES OF ANY KIND, EITHER EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO NON-INFRINGEMENT, MERCHANTABILITY OR FIT FOR A PARTICULAR
// PURPOSE.
// See the Mulan PSL v1 for more details.
// Description: main function
// Author: zhangsong
// Create: 2019-01-18

// go base main package
package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/docker/docker/pkg/pidfile"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"
	"isula.org/authz/authz"
	"isula.org/authz/core"
)

const (
	debugFlag      = "debug"
	policyFileFlag = "policy-file"
)

var (
	version = "0.1.0"
	pidFile = "/run/authz.pid"
)

func main() {

	app := cli.NewApp()
	app.Name = "authz-broker"
	app.Usage = "Authorization plugin for isulad"
	app.Version = version

	app.Action = func(c *cli.Context) {

		// init logrus
		logrus.SetFormatter(&logrus.TextFormatter{})
		logrus.SetOutput(os.Stdout)
		logrus.SetLevel(logrus.DebugLevel)
		if c.GlobalBool(debugFlag) {
			logrus.SetLevel(logrus.DebugLevel)
		} else {
			logrus.SetLevel(logrus.InfoLevel)
		}

		// init authz pid file
		file, err := pidfile.New(pidFile)
		if err != nil {
			panic(fmt.Errorf("create new pid file '%s' error pid file found, ensure authz is not running or delete %s", pidFile, pidFile))
		}
		if err := os.Chmod(pidFile, 0640); err != nil {
			panic(err)
		}
		defer func() {
			if err := file.Remove(); err != nil {
				panic(fmt.Errorf("remove pid file failed: %v", err))
			}
		}()

		// start authz server
		authorizer := authz.NewAuthorizer(c.GlobalString(policyFileFlag))
		auditor := authz.NewAuditor()
		srv := core.NewAuthZServer(authorizer, auditor)
		go func() {
			err = srv.Start()
			if err != nil {
				panic(err)
			}
		}()
		ch := make(chan os.Signal, 1)
		signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
		<-ch
		srv.Stop()
	}

	app.Flags = []cli.Flag{
		cli.BoolFlag{
			Name:   debugFlag,
			Usage:  "Enable debug mode",
			EnvVar: "DEBUG",
		},
		cli.StringFlag{
			Name:   policyFileFlag,
			Value:  "/var/lib/authz-broker/policy.json",
			EnvVar: "AUTHZ-POLICY-FILE",
			Usage:  "Specify authz policy file",
		},
	}

	app.Run(os.Args)
}
