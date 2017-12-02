GOPATH  = $(shell pwd)
PKG_DIR = src/tpi-mon

DB_PASSWORD := tpimon_dev123

.PHONY: all clean

all: local cloud mock

pkg:
	docker build -t tpi-mon-pkg -f $(PKG_DIR)/pkg/Dockerfile $(PKG_DIR)

local: pkg
	docker build -t tpi-mon-local -f $(PKG_DIR)/local/Dockerfile $(PKG_DIR)

cloud: pkg db
	docker build -t tpi-mon-cloud -f $(PKG_DIR)/cloud/Dockerfile $(PKG_DIR)

db:
	docker build --build-arg=DB_PASSWORD=$(DB_PASSWORD) -t tpi-mon-db db

mock: pkg
	docker build -t tpi-mon-mock -f $(PKG_DIR)/mock/Dockerfile $(PKG_DIR)

test: test-db

test-db:
	cd $(PKG_DIR)/cloud/db && go test
