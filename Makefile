.PHONY: all
all:
	$(MAKE) test

.PHONY: test
test:
	go test -v ./.
