go-dev:
	go run src/main.go

bun-dev:
	bun run dev ./src/client

bun-build:
	cd src/client ; bun run build ; cd ../..

bun-lint:
	cd src/client ; bun run lint ; cd ../..

dev:
	@make -j1 bun-dev go-dev

dev-build:
	@make -j1 bun-build go-dev