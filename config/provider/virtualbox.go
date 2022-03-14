// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

type VirtualBox struct {
	CheckGuestAdditions bool       `hcl:"check_guest_additions,optional"`
	CPUs                int32      `hcl:"cpus,optional"`
	Customize           [][]string `hcl:"customize,optional"`
	DefaultNICType      string     `hcl:"default_nic_type,optional"`
	GUI                 bool       `hcl:"gui,optional"`
	LinkedClone         bool       `hcl:"linked_clone,optional"`
	Memory              int32      `hcl:"memory,optional"`
	Name                string     `hcl:"name,optional"`
}
