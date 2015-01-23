pwd := $(shell pwd)
srcs := $(pwd)/../cmd
bins := solve show check

cases := 002_020_span
cases += 002_020_temp 002_040_temp 004_080_temp

all: $(bins)

define define_bin
$(1): $(shell find $(srcs)/$(1) -name '*.go')
	cd $(srcs)/$(1) && go build
	mv $(srcs)/$(1)/$(1) $(1)
endef

$(foreach bin,$(bins),$(eval $(call define_bin,$(bin))))

$(cases): $(bins)
	@./solve -c assets/$@.json -o solve_$@.mat
	@./show -c assets/$@.json -i solve_$@.mat
	@./check -c assets/$@.json -i solve_$@.mat

clean:
	rm -f $(bins) *.mat

.PHONY: all clean
