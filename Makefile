SOURCES=$(shell find ./src/acr/framework -name '*.go' -printf '%T@\t%p\n' | sort -k 1nr | cut -f2-)

GO=./pkg/linux_amd64/acr/framework.a
JS=./pkg/linux_js/acr/framework.a
MOBILE=./framework.aar

all: $(GO) $(JS) $(MOBILE)

$(GO): $(SOURCES)
	. ./vars.sh && go install acr/framework

$(JS): $(SOURCES)
	. ./vars.sh && gopherjs install acr/framework

$(MOBILE): $(SOURCES)
	. ./vars.sh && gomobile bind -target=android acr/framework
