
PKG_DIR = src/tpi-mon

.PHONY: all clean

all: tpi-local tpi-cloud tpi-mock

tpi-pkg:
	docker build -t tpi-mon-pkg -f $(PKG_DIR)/pkg/Dockerfile $(PKG_DIR)

tpi-local: tpi-pkg
	docker build -t tpi-mon-local -f $(PKG_DIR)/local/Dockerfile $(PKG_DIR)

tpi-cloud: tpi-pkg
	docker build -t tpi-mon-cloud -f $(PKG_DIR)/cloud/Dockerfile $(PKG_DIR)

tpi-mock: tpi-pkg
	docker build -t tpi-mon-mock -f $(PKG_DIR)/mock/Dockerfile $(PKG_DIR)
