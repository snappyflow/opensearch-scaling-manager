PLATFORM ?= linux
INCLUDESIM ?=false

export GO_BUILD=env go build
export SCALING_MANAGER_TAR_GZ="scaling_manager.tar.gz"
export SCALING_MANAGER_LIB="scaling_manager_lib"

default: build

#build: check clean init gotest
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

gotest:
	go test ./...

clean:
	go clean --cache
	mv scaling_manager.service sm.service
	rm -rf scaling_manager* go.mod go.sum
	mv sm.service scaling_manager.service

cleaninstall:
	systemctl stop scaling_manager
	systemctl disable scaling_manager
	rm -rf /usr/local/$(SCALING_MANAGER_LIB)
	rm -f /etc/systemd/system/scaling_manager.service

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
	cp config.yaml mappings.json scaling_manager.service $(SCALING_MANAGER_LIB)
    ifeq ($(INCLUDESIM),true)
	cp -rf simulator $(SCALING_MANAGER_LIB)
    endif
	cp logger/log_config.json $(SCALING_MANAGER_LIB)/logger
	cp provision/mappings.json $(SCALING_MANAGER_LIB)/provision
    ifeq ($(PLATFORM),windows)
		cp scaling_manager.exe $(SCALING_MANAGER_LIB)
    else
		cp scaling_manager $(SCALING_MANAGER_LIB)
    endif
	tar -czf $(SCALING_MANAGER_TAR_GZ) $(SCALING_MANAGER_LIB)

install: cleaninstall
	@if [ ! -f $(SCALING_MANAGER_TAR_GZ) ]; then \
		echo "The Scaling manager tarball is missing"; \
		exit 1; \
	fi
	rm -rf $(SCALING_MANAGER_LIB)
	tar -xzf $(SCALING_MANAGER_TAR_GZ)
    ifeq ($(PLATFORM),linux)
		tar -C /usr/local -xzf $(SCALING_MANAGER_TAR_GZ)
		chmod +x /usr/local/$(SCALING_MANAGER_LIB)/scaling_manager
		mv -f /usr/local/$(SCALING_MANAGER_LIB)/scaling_manager.service /etc/systemd/system/
		systemctl enable scaling_manager
		systemctl daemon-reload -l
    endif
