# GOPATH = $(shell pwd)/../..

# TARGETS = tpi-local tpi-mock tpi-cloud
# GO_FILES = $(shell find . -path ./cmd -prune -o -name '*.go')
# .PHONY: all clean

# all: $(TARGETS)


# PKG_DIRS = $(shell find ./pkg -name '*.go'  -exec dirname \{\} \; | uniq)
# PKG_FILES = $(shell find pkg -name '*.go')
# LOCAL_FILES = $(shell find local -name '*.go')
# CLOUD_FILES = $(shell find cloud -name '*.go')
# LOCAL_FILES = $(shell find local -name '*.go')
# MOCK_FILES = $(shell find mock -name '*.go')


# tpi-cloud: $(CLOUD_FILES) .pkg
# 	go build -o tpi-cloud ./cloud

# tpi-local: $(LOCAL_FILES) .pkg
# 	go build -o tpi-local ./local

# tpi-mock: $(MOCK_FILES) .pkg
# 	go build -o tpi-mock ./mock

# .pkg: $(PKG_FILES)
# 	go generate $(PKG_DIRS)
# 	go build $(PKG_DIRS)
# 	touch .pkg

# clean:
# 	rm -f $(TARGETS) .pkg

.PHONY: all clean

all: tpi-local tpi-cloud tpi-mock

tpi-pkg:
	docker build -t tpi-mon-pkg -f pkg/Dockerfile .

tpi-local: tpi-pkg
	docker build -t tpi-mon-local -f local/Dockerfile .

tpi-cloud: tpi-pkg
	docker build -t tpi-mon-cloud -f cloud/Dockerfile .

tpi-mock: tpi-pkg
	docker build -t tpi-mon-mock -f mock/Dockerfile .
