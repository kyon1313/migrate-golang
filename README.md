##### Golang Migrate Documentation ######

## Introduction
golang-migrate is a CLI tool and library for running database migrations. It supports various databases and allows you to manage schema changes in a controlled and versioned manner.


### Common Functions for golang-migrate
Here are some common functions you can use with golang-migrate:
#### Create a new migration:
    migrate create -ext sql -dir database/migrations -seq <migration_name>
#### Apply migrations:
    migrate -database ${DATABASE_URL} -path database/migrations up
#### Revert migrations:
    migrate -database ${DATABASE_URL} -path database/migrations down
#### Drop all tables:
    migrate -database ${DATABASE_URL} -path database/migrations drop


### Installation
To install golang-migrate, you can use Scoop:
    scoop install migrate
### If you don’t have Scoop installed, open PowerShell and run:
    Set-ExecutionPolicy RemoteSigned -Scope CurrentUser
    Invoke-RestMethod -Uri https://get.scoop.sh | Invoke-Expression


## Creating Migrations 
#### Create a directory for migrations
    mkdir database\migrations

### Using Makefile and Environment Variables
Create a .env file for your environment variables:

### Note on DATABASE_URL
#### When running it in Go, the DATABASE_URL should look like this:
    DATABASE_URL = sqlserver://$SA_USERNAME:$SA_PASSWORD@$DATABASE_HOST:$DATABASE_PORT?database=$DATABASE_NAME
#### When running it in the Makefile, the DATABASE_URL should look like this:
     DATABASE_URL = sqlserver://$(SA_USERNAME):$(SA_PASSWORD)@$(DATABASE_HOST):$(DATABASE_PORT)?database=$(DATABASE_NAME)


### Create a Makefile to Manage Migration Commands
#### migrate-create:
    migrate create -ext sql -dir database/migrations -seq $(name)
#### migrate-up:
    migrate -database ${DATABASE_URL} -path database/migrations up
#### migrate-down: 
    migrate -database ${DATABASE_URL} -path database/migrations down
#### migrate-version:
    migrate -database ${DATABASE_URL} -path database/migrations version


## Running Migrations Using the Makefile
#### Create a new migration:
    make create NAME=create_users_table
#### Apply migrations:
    make up
#### Revert migrations:
    make down
#### Check migration status:
    make version


## Go implementation will be demo with code , but heres how we can use it:
### Run migrations using Go
#### Create a new migration:
    go run main.go migrate-create <migration_name>
#### Apply migrations:
    go run main.go migrate-up
#### Revert migrations:
    go run main.go migrate-down
#### Check migration status:
    go run main.go migrate-version


For more detailed information, you can refer to the official:
- [golang-migrate documentation](https://github.com/golang-migrate/migrate)



## Running MSSQL Using Docker Compose ##### 
#### Build and run the Docker container:
    docker compose up --build
#### Access the MSSQL container:
    docker exec -it sql2 /bin/bash
#### Connect to MSSQL using sqlcmd:
    /opt/mssql-tools18/bin/sqlcmd -S tcp:localhost,1433 -U sa -P NewStrong@Passw0rd -C

## SQL COMMANDS FOR TESTING
#### Create a database:
    CREATE DATABASE sample;
    GO
#### List available databases:
    SELECT name FROM sys.databases;
    GO
#### Use the database:
    USE sample;
    GO






