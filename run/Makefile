bin := run
src := ../src
gos := $(wildcard $(src)/*.go)

cases := 002_020 002_040 004_080

all: $(bin)

$(bin): $(gos)
	cd $(src) && go build -o $@
	mv $(src)/$@ $@

$(cases): $(bin)
	./run solve -config assets/$@.json -output solve_$@.mat
	./run check -config assets/$@.json -input solve_$@.mat -output check_$@.mat

test:
	cd $(src) && go test

clean:
	rm -f $(bin) *.mat

.PHONY: all test clean
