PKGS=cmd \
  internal/handler \
  internal/domain \
  internal/repository \
  internal/entity

.PHONY: gomock
gomock: 
	go install github.com/vektra/mockery/v2@latest
	$(foreach dir, $(PKGS), \
		$(shell rm -rf $(dir)/gomocks) \
		$(foreach file, $(filter-out %_test.go, $(wildcard $(dir)/*.go)), \
			$(shell if [ $$(cat $(file) | grep -c "type .* interface") -gt 0 ]; \
			then \
				mockgen -source=$(file) -destination=$(dir)/gomocks/$(notdir $(file)) -package=gomocks; \
			fi ) \
		) \
	)

# Download dependencies
.PHONY: deps
deps:
	go mod download
	go mod tidy

# Run the application
.PHONY: run
run:
	go run cmd/main.go
