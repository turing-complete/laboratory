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

cases := 001_010_delay 001_010_energy 001_010_profile
cases += 004_040_delay 004_040_energy
cases += 009_090_delay 009_090_energy
cases += 016_160_delay 016_160_energy
cases += 025_250_delay 025_250_energy

define define_case
$(1)_config_base := assets/$(1).json
$(1)_hash_base := $$(shell hash md5sum 2>/dev/null && (cat $$($(1)_config_base) | md5sum | cut -d' ' -f1) || (cat $$($(1)_config_base) | md5))

$(1)_config := $$(shell test -e assets/$(1)_ext.json && echo assets/$(1)_ext.json || echo assets/$(1).json)
$(1)_hash := $$(shell hash md5sum 2>/dev/null && (cat $$($(1)_config) $$($(1)_config_base) | md5sum | cut -d' ' -f1) || (cat $$($(1)_config) $$($(1)_config_base) | md5))

$(1)_tgff := assets/$(shell echo $(1) | sed 's/_[a-z]*$$//').tgff

$(1)_common := $$($(1)_tgff)

$(1)_construct := data/$(1)_$$($(1)_hash)_construct.h5
$(1)_observe := data/$(1)_$$($(1)_hash_base)_observe.h5
$(1)_predict := data/$(1)_$$($(1)_hash)_predict.h5
$(1)_explore := data/$(1)_$$($(1)_hash)_explore.h5

$$($(1)_construct): construct $$($(1)_common)
	@./construct -c $$($(1)_config) -o $$@

$$($(1)_predict): $$($(1)_construct) predict $$($(1)_common)
	@./predict -c $$($(1)_config) -i $$< -o $$@

$$($(1)_observe): observe $$($(1)_common)
	@./observe -c $$($(1)_config) -o $$@

$(1): $$($(1)_predict) $$($(1)_observe) compare $$($(1)_common)
	@./compare -c $$($(1)_config) -i $$($(1)_predict) -o $$($(1)_observe)

$(1)_explore: $$($(1)_explore)

$$($(1)_explore): explore $$($(1)_common)
	@./explore -c $$($(1)_config) -o $$@ -v -n 50 -s "[1,2]"
endef

$(foreach case,$(cases),$(eval $(call define_case,$(case))))

assets/%:
	$(MAKE) -C assets $*

clean: flush
	rm -f $(bins)

flush:
	rm -f data/*.h5

.PHONY: all clean flush
