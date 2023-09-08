// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package network

import "github.com/hashicorp/hcl/v2"

type ForwardedPort struct {
	AutoCorrect bool   `hcl:"auto_correct,optional"`
	Guest       int32  `hcl:"guest,optional"`
	GuestIP     string `hcl:"guest_ip,optional"`
	Host        int32  `hcl:"host,optional"`
	HostIP      string `hcl:"host_ip,optional"`
	ID          string `hcl:"id,optional"`
	Protocol    string `hcl:"protocol,optional"`
	Type        string `hcl:"type,optional"`

	Body   hcl.Body `hcl:",body"`
	Remain hcl.Body `hcl:",remain"`
}

type Network struct {
	AutoConfig bool   `hcl:"auto_config,optional"`
	IP         string `hcl:"ip,optional"`
	Netmask    int32  `hcl:"netmask,optional"`

	Body   hcl.Body `hcl:",body"`
	Remain hcl.Body `hcl:",remain"`
}
