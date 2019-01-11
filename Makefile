.PHONY: test test-docker

test:
	go test -v

# run the integration tests in a docker container (linux platform target)
#
# this requires --privileged execution, but is still significantly faster than
# firing up a full Linux VM to test if developing on a different platform.
test-docker:
	docker run --rm --privileged \
		-e "GO111MODULE=on" \
		-v "$(CURDIR):/code" \
		--workdir "/code" \
		golang:1.11-alpine "go" "test" "-v"
