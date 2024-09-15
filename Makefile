ifeq (,$(wildcard .env))
    $(error .env file not found)
endif

include .env

create:
	go run main.go create $(NAME)

up:
	go run main.go up

down:
	go run main.go down

force:
	go run main.go force $(VERSION)

drop:
	go run main.go drop

version:
	go run main.go version
