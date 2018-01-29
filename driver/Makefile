.PHONY: image

NAME?=localflex
IMAGE?=localflex-driver

image: localflex
	docker build -t $(IMAGE) -f Dockerfile .

localflex: $(shell find . -name "main.go")
	GOOS=linux GOOARCH=amd64 CGO_ENABLE=0 go build -o $(NAME) .

clean:
	rm -f $(NAME)