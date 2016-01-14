root := $(dir $(shell pwd))

commands := sweep observe
commands += approximate predict
commands += compare

dependencies := $(shell find $(root)/src/internal -name '*.go')

all: $(commands)

define define_command
$(1): $(dependencies) $(shell find $(root)src/$(1) -name '*.go')
	cd $(root)src/$(1) && go build
	@mv $(root)src/$(1)/$(1) $(1)
endef

$(foreach command,$(commands),$(eval $(call define_command,$(command))))

cases := 001_010_delay 001_010_energy 001_010_temperature
cases += 001_020_delay 001_020_energy 001_020_temperature
cases += 001_030_delay 001_030_energy 001_030_temperature
cases += 001_040_delay 001_040_energy 001_040_temperature

cases += 004_010_delay 004_010_energy 004_010_temperature
cases += 004_020_delay 004_020_energy 004_020_temperature
cases += 004_030_delay 004_030_energy 004_030_temperature
cases += 004_040_delay 004_040_energy 004_040_temperature

cases += 008_010_delay 008_010_energy 008_010_temperature
cases += 008_020_delay 008_020_energy 008_020_temperature
cases += 008_030_delay 008_030_energy 008_030_temperature
cases += 008_040_delay 008_040_energy 008_040_temperature

define hash
$(1) := $$(shell hash md5sum 2>/dev/null && (cat $(2) | md5sum | cut -d' ' -f1) || (cat $(2) | md5))
endef

define define_case
$(1)_config_base := assets/$(1).json
$(call hash,$(1)_hash_base,assets/base.json $$($(1)_config_base))

$(1)_config := $$(shell test -e assets/$(1)_ext.json && echo assets/$(1)_ext.json || echo assets/$(1).json)
$(call hash,$(1)_hash,$$($(1)_config) $$($(1)_config_base))

$(1)_reference := data/$(1)_$$($(1)_hash_base)_reference.h5
$(1)_observe := data/$(1)_$$($(1)_hash_base)_observe.h5
$(1)_observe_sweep := data/$(1)_$$($(1)_hash_base)_observe_sweep.h5

$(1)_approximate := data/$(1)_$$($(1)_hash)_approximate.h5
$(1)_approximate_sweep := data/$(1)_$$($(1)_hash)_approximate_sweep.h5
$(1)_predict := data/$(1)_$$($(1)_hash)_predict.h5
$(1)_compare := data/$(1)_$$($(1)_hash)_compare.h5

$$($(1)_reference): observe
	@./observe -c $$($(1)_config) -s 0x42424242 -o $$@

$$($(1)_observe): observe
	@./observe -c $$($(1)_config) -o $$@

$$($(1)_observe_sweep): sweep
	@./sweep -c $$($(1)_config) -n 100 -s "[0, 1]" -o $$@

$$($(1)_approximate): approximate
	@./approximate -c $$($(1)_config) -o $$@

$$($(1)_approximate_sweep): $$($(1)_approximate) sweep
	@./sweep -c $$($(1)_config) -approximate $$< -n 100 -s "[0, 1]" -o $$@

$$($(1)_predict): $$($(1)_approximate) predict
	@./predict -c $$($(1)_config) -approximate $$< -o $$@

$$($(1)_compare): $$($(1)_reference) $$($(1)_observe) $$($(1)_predict) compare
	@./compare -c $$($(1)_config) -reference $$($(1)_reference) -observe $$($(1)_observe) -predict $$($(1)_predict) -o $$@

$(1)_reference: $$($(1)_reference)

$(1)_observe: $$($(1)_observe)

$(1)_observe_sweep: $$($(1)_observe_sweep)

$(1)_approximate: $$($(1)_approximate)

$(1)_approximate_sweep: $$($(1)_approximate_sweep)

$(1)_predict: $$($(1)_predict)

$(1)_compare: $$($(1)_compare)

$(1): $$($(1)_compare)
endef

$(foreach case,$(cases),$(eval $(call define_case,$(case))))

assets/%:
	$(MAKE) -C assets $*

clean: flush
	@rm -f $(commands)

flush:
	@rm -f data/*.h5

.DELETE_ON_ERROR:

.PHONY: all clean flush
