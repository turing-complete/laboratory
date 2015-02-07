pwd := $(shell pwd)
cmd := $(pwd)/../cmd
pkg := $(pwd)/../pkg

bins := construct sample compare
srcs := $(shell find $(cmd)/internal -name '*.go')
srcs += $(shell find $(pkg) -name '*.go')

all: $(bins)

define define_bin
$(1): $(srcs) $(shell find $(cmd)/$(1) -name '*.go')
	cd $(cmd)/$(1) && go build
	mv $(cmd)/$(1)/$(1) $(1)
endef

$(foreach bin,$(bins),$(eval $(call define_bin,$(bin))))

cases := 002_020_delay 004_040_delay 008_080_delay 016_160_delay 032_320_delay
cases += 002_020_energy 004_040_energy 008_080_energy 016_160_energy 032_320_energy
cases += 002_020_slice 004_040_slice 008_080_slice 016_160_slice 032_320_slice

define define_case
$(1)_config := assets/$(1).json
$(1)_hash := $$(shell hash md5sum 2>/dev/null && (cat $$($(1)_config) | md5sum | cut -d' ' -f1) || (cat $$($(1)_config) | md5))
$(1)_construct := $(1)_$$($(1)_hash)_construct.mat
$(1)_sample := $(1)_$$($(1)_hash)_sample.mat

$(1): $$($(1)_sample) compare
	@./compare -c $$($(1)_config) -i $$<

$$($(1)_sample): $$($(1)_construct) sample
	@./sample -c $$($(1)_config) -i $$< -o $$@

$$($(1)_construct): construct
	@./construct -c $$($(1)_config) -o $$@
endef

$(foreach case,$(cases),$(eval $(call define_case,$(case))))

clean:
	rm -f $(bins) *.mat

.PHONY: all clean
