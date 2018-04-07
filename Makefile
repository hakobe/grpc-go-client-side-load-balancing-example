.PHONY: build up down gen-pb
build:
	docker-compose build
up:
	docker-compose up
down:
	docker-compose down
gen-pb:
	mkdir -p trylb && \
	docker-compose run --rm --no-deps client /protoc/bin/protoc ./trylb.proto --go_out=plugins=grpc:trylb
