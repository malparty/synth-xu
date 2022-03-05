package racks

import (
	"github.com/malparty/synth-xu/lib/modules"
	"github.com/malparty/synth-xu/lib/modules/modulators"
)

type ModulesChain struct {
	moduleFuncs  []modules.ModuleFunction
	envelopeFunc modules.ModuleFunction
	Envelope     *modulators.Envelope
}

func NewChainFunc(envelope *modulators.Envelope, moduleItems []modules.Module) *ModulesChain {
	envFunc := envelope.GetModuleFunc()

	chain := &ModulesChain{
		Envelope:     envelope,
		moduleFuncs:  []modules.ModuleFunction{},
		envelopeFunc: envFunc,
	}

	for _, modules := range moduleItems {
		chain.moduleFuncs = append(chain.moduleFuncs, modules.GetModuleFunc())
	}

	return chain
}

func (m *ModulesChain) ChainFuncControlled(stat float64, delta float64) float64 {
	stat = m.ChainFunc(stat, delta)

	return m.envelopeFunc(stat, delta)
}

func (m *ModulesChain) ChainFunc(stat float64, delta float64) float64 {
	for _, moduleFunc := range m.moduleFuncs {
		stat = moduleFunc(stat, delta)
	}

	return stat
}
