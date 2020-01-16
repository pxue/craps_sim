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
	@(go run cmd/craps/main.go)
