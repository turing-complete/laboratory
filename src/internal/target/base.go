package target

import (
	"fmt"

	"github.com/turing-complete/laboratory/src/internal/config"
	"github.com/turing-complete/laboratory/src/internal/system"
	"github.com/turing-complete/laboratory/src/internal/uncertainty"
)

type base struct {
	uncertainty.Transform

	system *system.System
	config *config.Target

	ni uint
	no uint
}

func newBase(system *system.System, transform uncertainty.Transform,
	config *config.Target, ni, no uint) (base, error) {

	return base{
		Transform: transform,

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
