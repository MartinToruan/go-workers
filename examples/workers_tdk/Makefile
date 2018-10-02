# go build command
gobuild:
	@go build -v -o workers_tdk cmd/workers_tdk/*.go

# go run command
gorun:
	make gobuild
	@./workers_tdk --config_path='./config/tkp-app.{TKPENV}.yaml'

# deploy command
deploy:
	@echo "deploying"