package plugincore

import (
	"github.com/hashicorp/go-argmapper"
	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-plugin"

	"github.com/hashicorp/vagrant-plugin-sdk/core"
	"github.com/hashicorp/vagrant-plugin-sdk/internal-shared/protomappers"
	coreplugin "github.com/hashicorp/vagrant-plugin-sdk/internal/plugin/core"
)

func NewTargetPlugin(impl core.Target, log hclog.Logger) plugin.GRPCPlugin {
	return &coreplugin.TargetPlugin{
		Impl:    impl,
		Mappers: coreMappers,
		Logger:  log,
	}
}

func NewProjectPlugin(impl core.Project, log hclog.Logger) plugin.GRPCPlugin {
	return &coreplugin.ProjectPlugin{
		Impl:    impl,
		Mappers: coreMappers,
		Logger:  log,
	}
}

func NewBasisPlugin(impl core.Basis, log hclog.Logger) plugin.GRPCPlugin {
	return &coreplugin.BasisPlugin{
		Impl:    impl,
		Mappers: coreMappers,
		Logger:  log,
	}
}

var (
	coreMappers = []*argmapper.Func{}
)

func init() {
	for _, f := range coreplugin.MapperFns {
		fn, err := argmapper.NewFunc(f)
		if err != nil {
			panic(err)
		}

		coreMappers = append(coreMappers, fn)
	}
	for _, f := range protomappers.All {
		fn, err := argmapper.NewFunc(f)
		if err != nil {
			panic(err)
		}

		coreMappers = append(coreMappers, fn)
	}
}
