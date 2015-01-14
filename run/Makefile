bin := run
src := ../src
gos := $(wildcard $(src)/*.go)

cases := 002_020 002_040 004_080

all: $(bin)

$(bin): $(gos)
	cd $(src) && go build -o $@
	mv $(src)/$@ $@

$(cases): $(bin)
	./run solve -c assets/$@.json -o solve_$@.mat
	./run check -c assets/$@.json -i solve_$@.mat -o check_$@.mat

test:
	cd $(src) && go test

clean:
	rm -f $(bin) *.mat

.PHONY: all test clean
