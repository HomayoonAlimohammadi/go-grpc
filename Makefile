.PHONY: clean gen serve client test coverage 

# Variables
COVERAGE_THRESHOLD := 78

gen:
	protoc --go_out=. --go_opt=paths=import --go-grpc_out=. --go-grpc_opt=paths=import ./proto/*.proto 

serve:
	go run main.go serve

client:
	go run main.go client 

test: 
	rm -rf tmp
	mkdir tmp
	go test -cover -race -v -coverprofile=c.out ./...
	rm -rf tmp

coverage: test
	$(eval totalCoverage := $(shell go tool cover -func=c.out | grep total | grep -Eo '[0-9]+\.[0-9]+'))
	$(info Current test coverage: $(totalCoverage)%.)
	@if [ $(shell echo "$(totalCoverage) >= $(COVERAGE_THRESHOLD)" | bc -l) -eq 1 ]; then echo PASSED.; else echo FAILED. Threshold is $(COVERAGE_THRESHOLD)%.; false; fi

clean:
	rm -rf pb
	rm -rf tmp