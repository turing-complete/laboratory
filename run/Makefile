pwd := $(shell pwd)
cmd := $(pwd)/../cmd
pkg := $(pwd)/../pkg

bins := solve show check
srcs := $(shell find $(cmd)/internal -name '*.go')
srcs += $(shell find $(pkg) -name '*.go')

all: $(bins)

define define_bin
$(1): $(srcs) $(shell find $(cmd)/$(1) -name '*.go')
	cd $(cmd)/$(1) && go build
	mv $(cmd)/$(1)/$(1) $(1)
endef

$(foreach bin,$(bins),$(eval $(call define_bin,$(bin))))

cases := 002_020_span 004_040_span 008_080_span 016_160_span 032_320_span

define define_case
$(1)_hash := $(shell hash md5sum 2>/dev/null && (cat assets/$(1).json | md5sum | cut -d' ' -f1) || (cat assets/$(1).json | md5))

$(1): $(1)_$$($(1)_hash).mat show check
	@./show -c assets/$(1).json -i $$<
	@./check -c assets/$(1).json -i $$<

$(1)_$$($(1)_hash).mat: solve
	@./solve -c assets/$(1).json -o $$@
endef

$(foreach case,$(cases),$(eval $(call define_case,$(case))))

clean:
	rm -f $(bins) *.mat

.PHONY: all clean
