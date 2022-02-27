package modules

type ModulesChain struct {
	moduleFuncs []ModuleFunction
}

func NewChainFunc(modules []Module) *ModulesChain {
	chain := &ModulesChain{
		moduleFuncs: []ModuleFunction{},
	}
	for _, modules := range modules {
		chain.moduleFuncs = append(chain.moduleFuncs, modules.GetModuleFunc())
	}

	return chain
}

func (m *ModulesChain) ChainFunc(stat float64, delta float64) float64 {
	for _, moduleFunc := range m.moduleFuncs {
		stat = moduleFunc(stat, delta)
	}

	return stat
}
