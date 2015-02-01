package internal

type Target interface {
	Evaluate([]float64, []float64, []uint64)
	InputsOutputs() (uint32, uint32)
}
