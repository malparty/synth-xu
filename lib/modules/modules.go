package modules

type Module interface {
	GetModuleFunc() ModuleFunction
}

type ModuleFunction func(stat float64, delta float64) float64
