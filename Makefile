# 1. API
tfplan=terraform.tfplan
test_env=$(PWD)/test/env
compiled=$(PWD)/linker

build:
	GOOS=linux go build ./cmd/linker

clean:
	rm -f linker deployments/deploy.zip \
	docker rm -f linker-dynamodb

plan: build
	export TF_VAR_TARGET_BINARY="$(compiled)" && \
		cd deployments && terraform plan -out="$(tfplan)"

apply:
	cd deployments && terraform apply "$(tfplan)" \
		&& rm -f "$(tfplan)"

deploy: apply

dev: build
	./scripts/run-local.sh "$(compiled)" "$(test_env)"

# 2. Web
web-dev:
	cd web && npm start
web-build:
	cd web && npm run build
web-test:
	cd web && npm test
web-deploy:
	cd web && npm run deploy
