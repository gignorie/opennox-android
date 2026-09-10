format:
	clang-format --verbose --style=file -i ./src/*.c ./src/*.h ./src/*/*.c ./src/*/*.h ./src/*/*/*.c ./src/*/*/*.h  ./src/*/*/*/*.c ./src/*/*/*/*.h

build-server:
	cd ./src; \
	go run ./internal/noxbuild server

build-client:
	cd ./src; \
	go run ./internal/noxbuild client client-hd

build-client-win:
	cd ./src; \
	go run ./internal/noxbuild --os=windows client client-hd

build-server-docker:
	GIT_SHA=$$(git rev-parse --short HEAD); \
	GIT_TAG=$$(git name-rev --tags --name-only $$GIT_SHA); \
	docker build -t ghcr.io/opennox/opennox:dev -f ./docker/Dockerfile_server --target=server --build-arg GIT_SHA=$$GIT_SHA  --build-arg GIT_TAG=$$GIT_TAG ./src

build-server-demo-docker:
	GIT_SHA=$$(git rev-parse --short HEAD); \
	GIT_TAG=$$(git name-rev --tags --name-only $$GIT_SHA); \
	docker build -t ghcr.io/opennox/opennox:dev-demo -f ./docker/Dockerfile_server --target=demo --build-arg GIT_SHA=$$GIT_SHA  --build-arg GIT_TAG=$$GIT_TAG ./src

build-client-docker:
	GIT_SHA=$$(git rev-parse --short HEAD); \
	GIT_TAG=$$(git name-rev --tags --name-only $$GIT_SHA); \
	docker build -t ghcr.io/opennox/opennox-client:dev -f ./docker/Dockerfile_client --build-arg GIT_SHA=$$GIT_SHA  --build-arg GIT_TAG=$$GIT_TAG ./src
	mkdir -p ./build
	ID=$$(docker create ghcr.io/opennox/opennox-client:dev) && \
	docker cp $$ID:/home/runner/opennox/opennox ./build/ && \
	docker cp $$ID:/home/runner/opennox/opennox-hd ./build/ && \
	docker rm -f $$ID
