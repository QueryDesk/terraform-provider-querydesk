default: testacc

# Run acceptance tests
.PHONY: testacc
testacc:
	TF_ACC=1 go test ./... -v $(TESTARGS) -timeout 120m

testcover:
	TF_ACC=1 go test ./internal/... -v $(TESTARGS) -timeout 120m -coverprofile=coverage.out
