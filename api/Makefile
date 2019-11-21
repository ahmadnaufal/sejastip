TIMESTAMP = $(shell date +'%Y%m%d%H%M%S')
PWD			 	= $(shell pwd)

# LOCAL SETUP FOR COMPILED BINARY
GOOS	 ?= linux
GOARCH  = amd64

export GO111MODULE ?= on

all: compile test

test:
	go test -cover -coverprofile=cover.out -v ./...

coverage:
	go tool cover -html=cover.out -o cover.html
	open cover.html

bin:
	go build -o sejastip-api app/main.go

mod:
	go mod tidy

compile:
	GOOS=$(GOOS) GOARCH=$(GOARCH) CGO_ENABLED=0 go build -o sejastip-api app/main.go

migration:
	@echo "class Migration < ActiveRecord::Migration[5.1]\n  def up\n  end\n\n  def down\n  end\nend" > db/migrate/$(TIMESTAMP)_migration.rb

setup:
	docker run --rm -it --network host -v "$(PWD)/db":/app/db -v "$(PWD)/.env":/app/.env ahmadnaufal/standalone-migration:1.0.0 db:create
	docker run --rm -it --network host -v "$(PWD)/db":/app/db -v "$(PWD)/.env":/app/.env ahmadnaufal/standalone-migration:1.0.0 db:migrate

migrate:
	docker run --rm -it --network host -v "$(PWD)/db":/app/db -v "$(PWD)/.env":/app/.env ahmadnaufal/standalone-migration:1.0.0 db:migrate

migrate-rollback:
	docker run --rm -it --network host -v "$(PWD)/db":/app/db -v "$(PWD)/.env":/app/.env ahmadnaufal/standalone-migration:1.0.0 db:rollback

seed:
	docker run --rm -it --network host -v "$(PWD)/db":/app/db -v "$(PWD)/.env":/app/.env ahmadnaufal/standalone-migration:1.0.0 db:seed
