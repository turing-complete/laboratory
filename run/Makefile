pwd := $(shell pwd)
cmd := $(pwd)/../cmd
pkg := $(pwd)/../pkg

bins := solve show check
srcs := $(shell find $(cmd)/internal -name '*.go')
srcs += $(shell find $(pkg) -name '*.go')

cases := 002_020_span 004_040_span 008_080_span 016_160_span 032_320_span

all: $(bins)

define define_bin
$(1): $(srcs) $(shell find $(cmd)/$(1) -name '*.go')
	cd $(cmd)/$(1) && go build
	mv $(cmd)/$(1)/$(1) $(1)
endef

$(foreach bin,$(bins),$(eval $(call define_bin,$(bin))))

$(cases): $(bins)
	@./solve -c assets/$@.json -o solve_$@.mat
	@./show -c assets/$@.json -i solve_$@.mat
	@./check -c assets/$@.json -i solve_$@.mat

clean:
	rm -f $(bins) *.mat

.PHONY: all clean
