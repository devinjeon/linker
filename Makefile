tfplan=terraform.tfplan

build: cmd/linker
	GOOS=linux go build ./cmd/linker
clean:
	rm -f linker deployments/deploy.zip \
	docker rm -f linker-dynamodb
plan: build tf-plan
apply: tf-apply clean
tf-plan: deployments
	export TF_VAR_TARGET_BINARY="$(PWD)/linker" && \
		cd deployments && terraform plan -out="$(tfplan)"
tf-apply: deployments/$(tfplan)
	cd deployments && terraform apply "$(tfplan)" \
		&& rm -f "$(tfplan)"

test: build
	./scripts/test.sh "$(PWD)/linker"

web-start: web
	cd web && npm start
web-build: web
	cd web && npm run build
web-test: web
	cd web && npm test
