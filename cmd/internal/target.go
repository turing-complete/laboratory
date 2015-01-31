package internal

type Target interface {
	InputsOutputs() (uint32, uint32)
	Evaluate([]float64, []float64, []uint64)
}
