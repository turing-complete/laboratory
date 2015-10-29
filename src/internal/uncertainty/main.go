package uncertainty

type Uncertainty interface {
	Len() int
	Transform([]float64) []float64
}
