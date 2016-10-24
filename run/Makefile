OUTPUT_DIR ?= results

root := $(dir $(shell pwd))

delay := 002_010_delay 002_020_delay
delay += 004_010_delay 004_020_delay
delay += 008_010_delay 008_020_delay
delay += 002_004_delay

energy := 002_010_energy 002_020_energy
energy += 004_010_energy 004_020_energy
energy += 008_010_energy 008_020_energy

temperature := 002_010_temperature 002_020_temperature
temperature += 004_010_temperature 004_020_temperature
temperature += 008_010_temperature 008_020_temperature

cases := $(delay) $(energy) $(temperature)

commands := observe sweep
commands += approximate predict
commands += compare

dependencies := $(shell find $(root)/src/internal -name '*.go')

build: $(commands)

solve: $(cases)

delay: $(delay)

energy: $(energy)

temperature: $(temperature)

test:
	cd $(root) && go test ./src/...

define define_command
$(1): $(dependencies) $(shell find $(root)src/$(1) -name '*.go')
	cd $(root)src/$(1) && go build
	@mv $(root)src/$(1)/$(1) $(1)
endef

$(foreach command,$(commands),$(eval $(call define_command,$(command))))

define hash
$(1) := $$(shell hash md5sum 2> /dev/null && \
	(cat $(2) | md5sum | cut -d' ' -f1) || (cat $(2) | md5))
endef

define define_case
$(1)_parts := $(subst _, ,$(1))
$(1)_cores := $$(word 1,$$($(1)_parts))
$(1)_quantity := $$(word 3,$$($(1)_parts))

$(1)_configs := assets/base.json
$(call hash,$(1)_hash_base,$$($(1)_configs))

$(1)_configs += assets/$$($(1)_quantity).json
$(1)_configs += assets/$$($(1)_cores)_$$($(1)_quantity).json
$(1)_configs += assets/$(1).json
$(call hash,$(1)_hash,$$($(1)_configs))

$(1)_reference := $(OUTPUT_DIR)/$(1)_$$($(1)_hash_base)_reference.h5
$(1)_observe := $(OUTPUT_DIR)/$(1)_$$($(1)_hash_base)_observe.h5
$(1)_observe_sweep := $(OUTPUT_DIR)/$(1)_$$($(1)_hash_base)_observe_sweep.h5

$(1)_approximate := $(OUTPUT_DIR)/$(1)_$$($(1)_hash)_approximate.h5
$(1)_approximate_sweep := $(OUTPUT_DIR)/$(1)_$$($(1)_hash)_approximate_sweep.h5
$(1)_predict := $(OUTPUT_DIR)/$(1)_$$($(1)_hash)_predict.h5
$(1)_compare := $(OUTPUT_DIR)/$(1)_$$($(1)_hash)_compare.h5

$(1)_config := assets/$(1).json

$$($(1)_reference): observe
	@./observe -c $$($(1)_config) -s 0x42424242 -o $$@

$$($(1)_observe): observe
	@./observe -c $$($(1)_config) -n 10000 -o $$@

$$($(1)_observe_sweep): sweep
	@./sweep -c $$($(1)_config) -n 100 -s "[0, 1]" -o $$@

$$($(1)_approximate): approximate
	@./approximate -c $$($(1)_config) -o $$@

$$($(1)_approximate_sweep): $$($(1)_approximate) sweep
	@./sweep -c $$($(1)_config) -approximate $$< -n 100 -s "[0, 1]" -o $$@

$$($(1)_predict): $$($(1)_approximate) predict
	@./predict -c $$($(1)_config) -approximate $$< -o $$@

$$($(1)_compare): $$($(1)_reference) $$($(1)_observe) $$($(1)_predict) compare
	@./compare -c $$($(1)_config) -reference $$($(1)_reference) \
		-observe $$($(1)_observe) -predict $$($(1)_predict) -o $$@

$(1)_header:
	@echo "----/--/-- --:--:-- Case: $(1)"

$(1)_footer:
	@echo "----/--/-- --:--:-- Done: $(1)"

$(1)_reference: $$($(1)_reference)

$(1)_observe: $$($(1)_observe)

$(1)_observe_sweep: $$($(1)_observe_sweep)

$(1)_approximate: $$($(1)_approximate)

$(1)_approximate_sweep: $$($(1)_approximate_sweep)

$(1)_predict: $$($(1)_predict)

$(1)_compare: $$($(1)_compare)

$(1): $(1)_header $$($(1)_compare) $(1)_footer
endef

$(foreach case,$(cases),$(eval $(call define_case,$(case))))

assets/%:
	$(MAKE) -C assets $*

clean: flush
	@rm -f $(commands)

flush:
	@rm -f $(OUTPUT_DIR)/*.h5

.DELETE_ON_ERROR:

.PHONY: build clean flush solve test
