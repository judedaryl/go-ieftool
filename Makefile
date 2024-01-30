ARCH := $(shell uname -m)
ifeq ($(ARCH),x86_64)
	ARCH := amd64
endif
OS := $(shell uname -s | tr '[:upper:]' '[:lower:]')
VERSION := $(shell curl -s https://api.github.com/repos/Schumann-IT/go-ieftool/releases/latest | grep "tag_name" | awk '{print $$2}' | sed 's|[\"\,]*||g')

ieftool:
	@curl -s -L -o ieftool https://github.com/Schumann-IT/go-ieftool/releases/download/$(VERSION)/ieftool-$(OS)-$(ARCH)
	@chmod +x ieftool
	@if [ "$(OS)" = "darwin" ]; then\
        xattr -d com.apple.quarantine ./ieftool /dev/null 2>&1 | true; \
    fi

install: ieftool
	@echo "Installing ieftool $(VERSION) for $(OS)/$(ARCH)"
	@sudo mv ieftool /usr/local/bin/ieftool

clean:
	@rm -Rf build
	@rm -f ./ieftool