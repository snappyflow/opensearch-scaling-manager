PLATFORM ?= linux

export GO_BUILD=env go build
export SCALING_MANAGER_TAR_GZ="scaling_manager.tar.gz"
export SCALING_MANAGER_LIB="scaling_manager_lib"

default: build

build: check clean init
	go vet
	@echo $(PLATFORM)
    ifeq ($(PLATFORM), linux)
	GOOS=linux GOARCH=amd64 $(GO_BUILD) -o scaling_manager
    else ifeq ($(PLATFORM), windows)
	GOOS=windows GOARCH=amd64 $(GO_BUILD) -o scaling_manager.exe
    else ifeq ($(PLATFORM), macintosh)
	GOOS=darwin GOARCH=amd64 $(GO_BUILD) -o scaling_manager
    else
	@echo "Please provide correct PLATFORM[linux, windows, macintosh]"
	exit 1
    endif

fmt:
	@echo "==> Formatting source code with gofmt..."
	gofmt -s -w .

clean:
	go clean --cache
	rm -rf scaling_manager* go.mod go.sum

init:
	go mod init scaling_manager
	go mod tidy

clobber:
	go clean --cache
	rm -rf scaling_manager*

check:
    ifneq ($(PLATFORM),linux)
    ifneq ($(PLATFORM),windows)
    ifneq ($(PLATFORM),macintosh)
    $(error "Please provide correct PLATFORM[linux, windows, macintosh]")
    endif
    endif
    endif

pack: check
	rm -rf $(SCALING_MANAGER_LIB) $(SCALING_MANAGER_TAR_GZ)
	mkdir -p $(SCALING_MANAGER_LIB)
	mkdir -p $(SCALING_MANAGER_LIB)/logger
	mkdir -p $(SCALING_MANAGER_LIB)/provision
	cp -rf config.yaml mappings.json simulator $(SCALING_MANAGER_LIB)
	cp logger/log_config.json $(SCALING_MANAGER_LIB)/logger
	cp provision/mappings.json $(SCALING_MANAGER_LIB)/provision
    ifeq ($(PLATFORM),windows)
		cp scaling_manager.exe $(SCALING_MANAGER_LIB)
    else
		cp scaling_manager $(SCALING_MANAGER_LIB)
    endif
	tar -czf $(SCALING_MANAGER_TAR_GZ) $(SCALING_MANAGER_LIB)

install:
	rm -rf $(SCALING_MANAGER_LIB)
	@if [ ! -f $(SCALING_MANAGER_TAR_GZ) ]; then \
		echo "The Scaling manager tarball is missing"; \
		exit 1; \
	fi
	tar -xzf $(SCALING_MANAGER_TAR_GZ)
