package core

import "github.com/hashicorp/vagrant-plugin-sdk/proto/vagrant_plugin_sdk"

type TargetIndex interface {
	Delete(entry Target) (err error)
	Get(ref *vagrant_plugin_sdk.Ref_Target) (entry Target, err error)
	Includes(ref *vagrant_plugin_sdk.Ref_Target) (exists bool, err error)
	Set(entry Target) (updatedEntry Target, err error)
	// Recover(entry Target) (updatedEntry Target, err error)
	All() (targets []Target, err error)
}
