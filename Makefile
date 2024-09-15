ifeq (,$(wildcard .env))
    $(error .env file not found)
endif

include .env


create:
	migrate create -ext sql -dir database/migrations -seq $(NAME) 

up:
	migrate -database $(DATABASE_URL) -path database/migrations up 1

down:
	migrate -database $(DATABASE_URL) -path database/migrations down 1

force:
	migrate -database $(DATABASE_URL) -path database/migrations force $(VERSION)

drop:
	migrate -database $(DATABASE_URL) -path database/migrations drop

version:
	migrate -database $(DATABASE_URL) -path database/migrations version

