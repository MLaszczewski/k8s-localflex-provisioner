.PHONY: image

IMAGE?=localflex-provisioner

image: localflex-provisioner
	docker build -t $(IMAGE) -f Dockerfile .

localflex-provisioner: $(shell find . -name "*.go")
	CGO_ENABLED=0 go build -a -ldflags '-extldflags "-static"' -o localflex-provisioner .

clean:
	rm -f localflex-provisioner
