// Copyright (c) Huawei Technologies Co., Ltd. 2019. All rights reserved.
// authz is licensed under the Mulan PSL v2.
// You can use this software according to the terms and conditions of the Mulan PSL v2.
// You may obtain a copy of Mulan PSL v2 at:
//    http://license.coscl.org.cn/MulanPSL2
// THIS SOFTWARE IS PROVIDED ON AN "AS IS" BASIS, WITHOUT WARRANTIES OF ANY KIND, EITHER EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO NON-INFRINGEMENT, MERCHANTABILITY OR FIT FOR A PARTICULAR
// PURPOSE.
// See the Mulan PSL v2 for more details.
// Description: convert a method/url pattern to corresponding isulad action,
//              available for isulad
// Author: zhangsong
// Create: 2019-01-18

package authz

import (
	"regexp"
	"strings"
)

type routeslice []route

type route struct {
	pattern string
	method  string
	action  string
}

// isulad routes
var isuladRoutes = []route{
	{pattern: "/events", method: "GET", action: "isulad_events"},
	{pattern: "/version", method: "GET", action: "isulad_version"},
	{pattern: "/auth", method: "POST", action: "isulad_auth"},
	{pattern: "/_ping", method: "GET", action: "isulad_ping"},
	{pattern: "/info", method: "GET", action: "isulad_info"},
}

// image routes
var imageRoutes = []route{
	{pattern: "/build", method: "POST", action: "image_build"},
	{pattern: "/images/.+/get", method: "GET", action: "images_archive"},
	{pattern: "/images/search", method: "GET", action: "images_search"},
	{pattern: "/images/.+/tag", method: "POST", action: "image_tag"},
	{pattern: "/images/.+/json", method: "GET", action: "image_inspect"},
	{pattern: "/images/.+", method: "DELETE", action: "image_delete"},
	{pattern: "/images/.+/history", method: "GET", action: "image_history"},
	{pattern: "/images/.+/push", method: "POST", action: "image_push"},
	{pattern: "/images/create", method: "POST", action: "image_create"},
	{pattern: "/images/load", method: "POST", action: "images_load"},
	{pattern: "/images/json", method: "GET", action: "image_list"},
}

// volume routes
var volumeRoutes = []route{
	{pattern: "/volumes/.+", method: "GET", action: "volume_inspect"},
	{pattern: "/volumes", method: "GET", action: "volume_list"},
	{pattern: "/volumes/create", method: "POST", action: "volume_create"},
	{pattern: "/volumes/.+", method: "DELETE", action: "volume_remove"},
}

// nework routes
var networkRoutes = []route{
	{pattern: "/networks/.+", method: "GET", action: "network_inspect"},
	{pattern: "/networks", method: "GET", action: "network_list"},
	{pattern: "/networks/create", method: "POST", action: "network_create"},
	{pattern: "/networks/.+/connect", method: "POST", action: "network_connect"},
	{pattern: "/networks/.+/disconnect", method: "POST", action: "network_disconnect"},
	{pattern: "/networks/.+", method: "DELETE", action: "network_remove"},
}

// container routes
var containerRoutes = []route{
	{pattern: "/commit", method: "POST", action: "container_commit"},
	{pattern: "/containers/.+/wait", method: "POST", action: "container_wait"},
	{pattern: "/containers/.+/resize", method: "POST", action: "container_resize"},
	{pattern: "/containers/.+/export", method: "GET", action: "container_export"},
	{pattern: "/containers/.+/stop", method: "POST", action: "container_stop"},
	{pattern: "/containers/.+/kill", method: "POST", action: "container_kill"},
	{pattern: "/containers/.+/restart", method: "POST", action: "container_restart"},
	{pattern: "/containers/.+/start", method: "POST", action: "container_start"},
	{pattern: "/containers/.+/update", method: "POST", action: "container_update"},
	{pattern: "/containers/.+/exec", method: "POST", action: "container_exec_create"},
	{pattern: "/containers/.+/unpause", method: "POST", action: "container_unpause"},
	{pattern: "/containers/.+/pause", method: "POST", action: "container_pause"},
	{pattern: "/containers/.+/copy", method: "POST", action: "container_copyfiles"},
	{pattern: "/containers/.+/archive", method: "PUT", action: "container_archive_extract"},
	{pattern: "/containers/.+/archive", method: "HEAD", action: "container_archive_info"},
	{pattern: "/containers/.+/archive", method: "GET", action: "container_archive"},
	{pattern: "/containers/.+/attach/ws", method: "GET", action: "container_attach_websocket"},
	{pattern: "/containers/.+/attach", method: "POST", action: "container_attach"},
	{pattern: "/containers/json", method: "GET", action: "container_list"},
	{pattern: "/containers/.+/json", method: "GET", action: "container_inspect"},
	{pattern: "/containers/.+", method: "DELETE", action: "container_delete"},
	{pattern: "/containers/.+/rename", method: "POST", action: "container_rename"},
	{pattern: "/containers/.+/stats", method: "GET", action: "container_stats"},
	{pattern: "/containers/.+/changes", method: "GET", action: "container_changes"},
	{pattern: "/containers/.+/top", method: "GET", action: "container_top"},
	{pattern: "/containers/.+/logs", method: "GET", action: "container_logs"},
	{pattern: "/containers/create", method: "POST", action: "container_create"},
	{pattern: "/exec/.+/json", method: "GET", action: "container_exec_inspect"},
	{pattern: "/exec/.+/start", method: "POST", action: "container_exec_start"},
}

var routes = []routeslice{
	isuladRoutes,
	imageRoutes,
	volumeRoutes,
	networkRoutes,
	containerRoutes,
}

// ParseRoute convert a method/url pattern to corresponding isulad action
func ParseRoute(method, url string) string {
	for _, rs := range routes {
		for _, route := range rs {
			if route.method == method {
				pattern := strings.Replace(route.pattern, ".+", "[a-zA-Z0-9_.:/-]+", 1) + "$"
				url = strings.Split(url, "?")[0]
				match, err := regexp.MatchString(pattern, url)
				if err == nil && match {
					return route.action
				}

			}
		}
	}
	return ""
}
