vendor:
	GO111MODULE=on go mod vendor
	git add vendor && git commit -s -m "dependencies: update vendored libs"

.PHONY: vendor
