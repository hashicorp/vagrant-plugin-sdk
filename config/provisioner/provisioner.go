// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provisioner

type File struct {
	Source      string `hcl:"source"`
	Destination string `hcl:"destination"`
}

type Shell struct {
	Inline string `hcl:"inline,optional"`
	Path   string `hcl:"path,optional"`

	Args                          []string          `hcl:"args,optional"`
	Binary                        bool              `hcl:"binary,optional"`
	Env                           map[string]string `hcl:"env,optional"`
	KeepColor                     bool              `hcl:"keep_color.optional"`
	MD5                           string            `hcl:"md5,optional"`
	Name                          string            `hcl:"name,optional"`
	PowershellArgs                []string          `hcl:"powershell_args,optional"`
	PowershellElevatedInteractive bool              `hcl:"powershell_elevated_interactive,optional"`
	Privileged                    bool              `hcl:"privileged,optional"`
	Reboot                        bool              `hcl:"reboot,optional"`
	Reset                         bool              `hcl:"reset,optional"`
	SHA1                          string            `hcl:"sha1,optional"`
	SHA256                        string            `hcl:"sha256,optional"`
	SHA384                        string            `hcl:"sha384,optional"`
	SHA512                        string            `hcl:"sha512,optional"`
	Sensitive                     bool              `hcl:"sensitive,optional"`
	UploadPath                    string            `hcl:"upload_path,optional"`
}
