package quantity

import (
	"fmt"

	"github.com/turing-complete/laboratory/src/internal/config"
	"github.com/turing-complete/laboratory/src/internal/system"
	"github.com/turing-complete/laboratory/src/internal/uncertainty"
)

type base struct {
	uncertainty.Uncertainty

	system *system.System
	config *config.Quantity

	ni uint
	no uint
}

func newBase(system *system.System, uncertainty uncertainty.Uncertainty,
	config *config.Quantity, ni, no uint) (base, error) {

	return base{
		Uncertainty: uncertainty,

		system: system,
		config: config,

		ni: ni,
		no: no,
	}, nil
}

func (self *base) Dimensions() (uint, uint) {
	return self.ni, self.no
}

func (self *base) String() string {
	return fmt.Sprintf(`{inputs:%d outputs:%d}`, self.ni, self.no)
}
