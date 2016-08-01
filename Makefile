all: image

sadbangbotd:
	CGO_ENABLED=0 go build -o deploy/sadbangbotd github.com/ralfonso/sadbangbot/cmd/sadbangbotd
.PHONY: sadbangbotd

image: sadbangbotd
	docker build -t ralfonso/sadbangbot .
.PHONY: image
