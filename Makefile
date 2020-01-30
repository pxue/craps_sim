.PHONY: help run

all:
	@echo "****************************"
	@echo "** Craps build tool       **"
	@echo "****************************"
	@echo "make <cmd>"
	@echo ""
	@echo "commands:"
	@echo "  run                   - run in dev mode"
	@echo ""
	@echo ""

print-%: ; @echo $*=$($*)

run:
	@read -r -p "output filename: " fname;             \
	go run cmd/craps/main.go -file $${fname}

run-debug:
	@(go run cmd/craps/main.go -iter 1 -debug)
