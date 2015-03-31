pwd := $(shell pwd)
cmd := $(dir $(pwd))cmd
pkg := $(dir $(pwd))pkg

bins := construct predict observe compare explore
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
cases += 002_020_slice
cases += 002_020_switch
cases += 002_020_profile 004_040_profile

define define_case
$(1)_config_base := assets/$(1).json
$(1)_hash_base := $$(shell hash md5sum 2>/dev/null && (cat $$($(1)_config_base) | md5sum | cut -d' ' -f1) || (cat $$($(1)_config_base) | md5))

$(1)_config := $$(shell test -e assets/$(1)_ext.json && echo assets/$(1)_ext.json || echo assets/$(1).json)
$(1)_hash := $$(shell hash md5sum 2>/dev/null && (cat $$($(1)_config) $$($(1)_config_base) | md5sum | cut -d' ' -f1) || (cat $$($(1)_config) $$($(1)_config_base) | md5))

$(1)_construct := data/$(1)_$$($(1)_hash)_construct.h5
$(1)_observe := data/$(1)_$$($(1)_hash_base)_observe.h5
$(1)_predict := data/$(1)_$$($(1)_hash)_predict.h5
$(1)_explore := data/$(1)_$$($(1)_hash)_explore.h5

$$($(1)_construct): construct
	@./construct -c $$($(1)_config) -o $$@

$$($(1)_predict): $$($(1)_construct) predict
	@./predict -c $$($(1)_config) -i $$< -o $$@

$$($(1)_observe): observe
	@./observe -c $$($(1)_config) -o $$@

$(1): $$($(1)_predict) $$($(1)_observe) compare
	@./compare -c $$($(1)_config) -i $$($(1)_predict) -o $$($(1)_observe)

$(1)_explore: $$($(1)_explore)

$$($(1)_explore): explore
	@./explore -c $$($(1)_config) -o $$@ -v
endef

$(foreach case,$(cases),$(eval $(call define_case,$(case))))

clean: flush
	rm -f $(bins)

flush:
	rm -f data/*.h5

.PHONY: all clean flush
