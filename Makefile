# go run -ldflags "-X github.com/nshttpd/oci-api-keygen/cmd.version=1.0.1-BETA -X github.com/nshttpd/oci-api-keygen/cmd.shortSha=`git rev-parse HEAD`" main.go version

VERSION=`cat VERSION`
SHORTSHA=`git rev-parse --short HEAD`

LDFLAGS=-X github.com/nshttpd/oci-api-keygen/cmd.version=$(VERSION)
LDFLAGS+=-X github.com/nshttpd/oci-api-keygen/cmd.shortSha=$(SHORTSHA)

build:
	go build -ldflags "$(LDFLAGS)" .

utils:
	go get github.com/mitchellh/gox
	go get github.com/tcnksm/ghr

deploy: utils
	pwd
	ls
	gox -parallel=4 -ldflags "$(LDFLAGS)" -output "dist/oci-api-keygen_{{.OS}}_{{.Arch}}"
	ghr -t $(GITHUB_TOKEN) -u $(CIRCLE_PROJECT_USERNAME) -r $(CIRCLE_PROJECT_REPONAME) --replace $(VERSION) dist/
