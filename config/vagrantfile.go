// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package config

import (
	"github.com/hashicorp/hcl/v2"
)

type Vagrantfile struct {
	Communicators []*Communicator `hcl:"communicator,block" json:",omitempty"`
	DefinedVms    []*VM           `hcl:"define_vm,block" json:",omitempty"`
	SSH           *SSH            `hcl:"ssh,block" json:",omitempty"`
	Vagrant       *Vagrant        `hcl:"vagrant,block" json:",omitempty"`
	VM            *VM             `hcl:"vm,block" json:",omitempty"`

	// These are values which are set after finalizations

	DefinedVmKeys []string `json:",omitempty"`
	ListVMs       []*VM    `json:",omitempty"`

	Body   hcl.Body `hcl:",body" json:"-"`
	Remain hcl.Body `hcl:",remain" json:"-"`
}

type SSH struct {
	ConnectTimeout *int32 `hcl:"connect_timeout,optional" json:",omitempty"`
}

type Communicator struct {
	Name string `hcl:"name,label"`

	Body   hcl.Body `hcl:",body" json:"-"`
	Remain hcl.Body `hcl:",remain" json:"-"`
}

type Vagrant struct {
	Sensitive []string `hcl:"sensitive,optional" json:",omitempty"`
	Host      *string  `hcl:"host,optional" json:",omitempty"`
	Plugins   []Plugin `hcl:"plugins,block" json:",omitempty"`

	Body   hcl.Body `hcl:",body" json:"-"`
	Remain hcl.Body `hcl:",remain" json:"-"`
}

type Target struct {
	Name string `hcl:"name,label" json:"name"`

	Providers []*Provider `hcl:"provider,optional" json:",omitempty"`

	Provider *Provider `json:",omitempty"`
}

type VM struct {
	Name string `hcl:"name,label" json:"name"`

	AllowFstabModification     *bool             `hcl:"allow_fstab_modification,optional" json:",omitempty"`
	AllowHostsModification     *bool             `hcl:"allow_hosts_modification,optional" json:",omitempty"`
	AutoStart                  *bool             `hcl:"autostart,optional" json:",omitempty"`
	BaseMAC                    *string           `hcl:"base_mac,optional" json:",omitempty"`
	BootTimeout                *int32            `hcl:"boot_timeout,optional" json:",omitempty"`
	Box                        *string           `hcl:"box,optional" json:",omitempty"`
	BoxCheckUpdate             *bool             `hcl:"box_check_update,optional" json:",omitempty"`
	BoxDownloadChecksum        *string           `hcl:"box_download_checksum,optional" json:",omitempty"`
	BoxDownloadChecksumType    *string           `hcl:"box_download_checksum_type,optional" json:",omitempty"`
	BoxDownloadClientCert      *string           `hcl:"box_download_client_cert,optional" json:",omitempty"`
	BoxDownloadCACert          *string           `hcl:"box_download_ca_cert,optional" json:",omitempty"`
	BoxDownloadCAPath          *string           `hcl:"box_download_ca_path,optional" json:",omitempty"`
	BoxDownloadOptions         map[string]string `hcl:"box_download_options,optional" json:",omitempty"`
	BoxDownloadInsecure        *bool             `hcl:"box_download_insecure,optional" json:",omitempty"`
	BoxDownloadLocationTrusted *bool             `hcl:"box_download_location_trusted,optional" json:",omitempty"`
	BoxURL                     *string           `hcl:"box_url,optional" json:",omitempty"`
	BoxVersion                 *string           `hcl:"box_version,optional" json:",omitempty"`
	//	CloudInit                  *CloudInit        `hcl:"cloud_init,block" json:",omitempty"`
	Communicator *string `hcl:"communicator,optional" json:",omitempty"`
	// Disk                       *Disk             `hcl:"disk,optional" json:",omitempty"`
	GracefulHaltTimeout  *int32  `hcl:"graceful_halt_timeout,optional" json:",omitempty"`
	Guest                *string `hcl:"guest,optional" json:",omitempty"`
	Hostname             *string `hcl:"hostname,optional" json:",omitempty"`
	IgnoreBoxVagrantfile *bool   `hcl:"ignore_box_vagrantfile,optional" json:",omitempty"`
	// Networks                   []*Network        `hcl:"network,block" json:",omitempty"`
	PostUpMessage *string `hcl:"post_up_message,optional" json:",omitempty"`
	Primary       *bool   `hcl:"primary,optional" json:",omitempty"`
	// Providers                  []*Provider       `hcl:"provider,block" json:",omitempty"`
	// Provisioners               []*Provisioner    `hcl:"provisioner,block" json:",omitempty"`
	// SyncedFolders              []*SyncedFolder   `hcl:"synced_folder,block" json:",omitempty"`
	// UsablePortRange            *Range            `hcl:"usable_port_range,block" json:",omitempty"`

	//	Provider *Provider `json:",omitempty"`

	// Body   hcl.Body `hcl:",body" json:"-"`
	// Remain hcl.Body `hcl:",remain" json:"-"`
}

// Here with looking at DefinedVMs which are nested Vagrantfiles and how the
// naming will work with named subvms and the default. Pretty sure we are going
// to want to modify the structure on the Ruby side before we serialize and ship
// back over the wire

func (v *VM) Target() *Target {
	return &Target{
		Name: v.Name,
		// Provider:  v.Provider,
		// Providers: v.Providers,
	}
}

type CloudInit struct {
	ContentType *string `hcl:"content_type" json:",omitempty"`
	Inline      *string `hcl:"inline,optional" json:",omitempty"`
	Path        *string `hcl:"path,optional" json:",omitempty"`

	Body   hcl.Body `hcl:",body" json:"-"`
	Remain hcl.Body `hcl:",remain" json:"-"`
}

type Disk struct {
	Name string `hcl:"name,label" json:"name"`

	DiskExt        *string           `hcl:"disk_ext,optional" json:",omitempty"`
	File           *string           `hcl:"file,optional" json:",omitempty"`
	Primary        *bool             `hcl:"primary,optional" json:",omitempty"`
	ProviderConfig map[string]string `hcl:"provider_config,optional" json:",omitempty"`
	Size           *string           `hcl:"size"`

	Body   hcl.Body `hcl:",body" json:"-"`
	Remain hcl.Body `hcl:",remain" json:"-"`
}

type Network struct {
	Kind *string `hcl:"kind,label" json:"kind"`

	Body   hcl.Body `hcl:",body" json:"-"`
	Remain hcl.Body `hcl:",remain" json:"-"`
}

type Plugin struct {
	Name string `hcl:"name,label" json:"name"`

	EntryPoint *string  `hcl:"entry_point,optional" json:",omitempty"`
	Sources    []string `hcl:"sources,optional" json:",omitempty"`
	Version    *string  `hcl:"version,optional" json:",omitempty"`

	Body   hcl.Body `hcl:",body" json:"-"`
	Remain hcl.Body `hcl:",remain" json:"-"`
}

type Provider struct {
	Name string `hcl:"name,label" json:"name"`

	Body   hcl.Body `hcl:",body" json:"-"`
	Remain hcl.Body `hcl:",remain" json:"-"`
}

type Provisioner struct {
	Kind string `hcl:"kind,label" json:"kind"`

	Body   hcl.Body `hcl:",body" json:"-"`
	Remain hcl.Body `hcl:",remain" json:"-"`
}

type SyncedFolder struct {
	Source      string `hcl:"source,label" json:"source"`
	Destination string `hcl:"destination,label" json:"destination"`

	Disabled     *bool    `hcl:"disabled,optional" json:",omitempty"`
	Group        *string  `hcl:"group,optional" json:",omitempty"`
	ID           *string  `hcl:"id,optional" json:",omitempty"`
	MountOptions []string `hcl:"mount_options,optional" json:",omitempty"`
	Owner        *string  `hcl:"owner,optional" json:",omitempty"`
	Type         *string  `hcl:"type,optional" json:",omitempty"`

	Body   hcl.Body `hcl:",body" json:"-"`
	Remain hcl.Body `hcl:",remain" json:"-"`
}

// Type helpers

// Note: used for types.Range
type Range struct {
	Start int32 `hcl:"start" json:"start"`
	End   int32 `hcl:"end" json:"end"`
}
