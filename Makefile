.PHONY: all distclean prepare

OUTDIR=$(shell realpath release)

PROJ=exec-agent
VERSION=1.0.10
TIMESTAMP=`date +%s`

MAJOR=`echo $(VERSION)|cut -d'.' -f1`
MINOR=`echo $(VERSION)|cut -d'.' -f2`
PATCH=`echo $(VERSION)|cut -d'.' -f3`

BRANCH=`git rev-parse --abbrev-ref HEAD`
HASH=`git log -n1 --pretty=format:%h`
REVERSION=`git log --oneline|wc -l|tr -d ' '`
BUILD_TIME=`date +'%Y-%m-%d %H:%M:%S'`
LDFLAGS="-X 'main.gitBranch=$(BRANCH)' \
-X 'main.gitHash=$(HASH)' \
-X 'main.gitReversion=$(REVERSION)' \
-X 'main.buildTime=$(BUILD_TIME)' \
-X 'main.version=$(VERSION)'"

all: distclean linux.amd64 linux.386 windows.amd64 windows.386
	cp conf/manifest.yaml $(OUTDIR)/$(VERSION)/manifest.yaml
	cp CHANGELOG.md $(OUTDIR)/CHANGELOG.md
	rm -fr $(OUTDIR)/$(VERSION)/etc $(OUTDIR)/$(VERSION)/opt
version:
	@echo $(VERSION)
linux.amd64: prepare
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -ldflags $(LDFLAGS) \
		-o $(OUTDIR)/$(VERSION)/opt/$(PROJ)/bin/$(PROJ) main.go
	cd $(OUTDIR)/$(VERSION) && fakeroot tar -czvf $(PROJ)_$(VERSION)_linux_amd64.tar.gz \
		--warning=no-file-changed opt
	go run contrib/pack/release.go -o $(OUTDIR)/$(VERSION) \
		-conf contrib/pack/amd64.yaml \
		-name $(PROJ) -version $(VERSION) \
		-workdir $(OUTDIR)/$(VERSION) \
		-epoch $(REVERSION)
linux.386: prepare
	GOOS=linux GOARCH=386 CGO_ENABLED=0 go build -ldflags $(LDFLAGS) \
		-o $(OUTDIR)/$(VERSION)/opt/$(PROJ)/bin/$(PROJ) main.go
	cd $(OUTDIR)/$(VERSION) && fakeroot tar -czvf $(PROJ)_$(VERSION)_linux_386.tar.gz \
		--warning=no-file-changed opt
	go run contrib/pack/release.go -o $(OUTDIR)/$(VERSION) \
		-conf contrib/pack/i386.yaml \
		-name $(PROJ) -version $(VERSION) \
		-workdir $(OUTDIR)/$(VERSION) \
		-epoch $(REVERSION)
windows.amd64: prepare
	GOOS=windows GOARCH=amd64 CGO_ENABLED=0 go build -ldflags $(LDFLAGS) \
		-o $(OUTDIR)/$(VERSION)/opt/$(PROJ)/bin/$(PROJ).exe
	unix2dos conf/agent.conf
	makensis -DARCH=amd64 \
		-DPRODUCT_VERSION=$(VERSION) \
		-DBINDIR=$(OUTDIR)/$(VERSION)/opt/$(PROJ)/bin/$(PROJ).exe \
		-INPUTCHARSET UTF8 contrib/win.nsi
	mv contrib/$(PROJ)_$(VERSION)_windows_amd64.exe $(OUTDIR)/$(VERSION)/$(PROJ)_$(VERSION)_windows_amd64.exe
windows.386: prepare
	GOOS=windows GOARCH=386 CGO_ENABLED=0 go build -ldflags $(LDFLAGS) \
		-o $(OUTDIR)/$(VERSION)/opt/$(PROJ)/bin/$(PROJ).exe main.go
	unix2dos conf/agent.conf
	makensis -DARCH=386 \
		-DPRODUCT_VERSION=$(VERSION) \
		-DBINDIR=$(OUTDIR)/$(VERSION)/opt/$(PROJ)/bin/$(PROJ).exe \
		-INPUTCHARSET UTF8 contrib/win.nsi
	mv contrib/$(PROJ)_$(VERSION)_windows_386.exe $(OUTDIR)/$(VERSION)/$(PROJ)_$(VERSION)_windows_386.exe
prepare:
	rm -fr $(OUTDIR)/$(VERSION)/opt $(OUTDIR)/$(VERSION)/etc
	mkdir -p $(OUTDIR)/$(VERSION)/opt/$(PROJ)/bin \
		$(OUTDIR)/$(VERSION)/opt/$(PROJ)/conf
	cp conf/agent.conf $(OUTDIR)/$(VERSION)/opt/$(PROJ)/conf/agent.conf
	echo $(VERSION) > $(OUTDIR)/$(VERSION)/opt/$(PROJ)/.version
	sed -i "s|#MAJOR|$(MAJOR)|g" contrib/versioninfo.json
	sed -i "s|#MINOR|$(MINOR)|g" contrib/versioninfo.json
	sed -i "s|#PATCH|$(PATCH)|g" contrib/versioninfo.json
	sed -i "s|#VERSION|v$(VERSION)|g" contrib/versioninfo.json
distclean:
	rm -fr $(OUTDIR)