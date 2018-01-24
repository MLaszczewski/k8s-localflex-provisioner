.PHONY: localflex

NAME?=localflex

localflex: $(shell find . -name "main.go")
	GOOS=linux GOOARCH=amd64 CGO_ENABLE=0 go build -o $(NAME) .

clean:
	rm -f $(NAME)