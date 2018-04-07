.PHONY: build up down gen-pb
build:
	docker-compose build
up:
	docker-compose up
down:
	docker-compose down
gen-pb:
	cp guruguru.proto ./go/guruguru.proto && \
		docker-compose run --rm --no-deps go /protoc/bin/protoc ./try_lb.proto --go_out=plugins=grpc:guruguru
