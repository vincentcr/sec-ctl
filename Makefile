GOPATH  = $(shell pwd)
PKG_DIR = src/sec-ctl

DB_PASSWORD := secctl_dev123

.PHONY: all clean db

all: local cloud mock

pkg:
	docker build -t sec-ctl-pkg -f $(PKG_DIR)/pkg/Dockerfile $(PKG_DIR)

local: pkg
	docker build -t sec-ctl-local -f $(PKG_DIR)/local/Dockerfile $(PKG_DIR)

cloud: pkg db
	docker build -t sec-ctl-cloud -f $(PKG_DIR)/cloud/Dockerfile $(PKG_DIR)

db:
	docker build --build-arg=DB_PASSWORD=$(DB_PASSWORD) -t sec-ctl-db db

mock: pkg
	docker build -t sec-ctl-mock -f $(PKG_DIR)/mock/Dockerfile $(PKG_DIR)

test: test-cloud

test-cloud:
	docker-compose up -d test-cloud-db
	cd $(PKG_DIR)/cloud/db && go test
