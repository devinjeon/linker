tfplan=terraform.tfplan

build: cmd/linker
	GOOS=linux go build ./cmd/linker
clean:
	rm -f linker deployments/deploy.zip
plan: build tf-plan
apply: tf-apply clean
tf-plan: deployments
	export TF_VAR_TARGET_BINARY="$(PWD)/linker" && \
		cd deployments && terraform plan -out="$(tfplan)"
tf-apply: deployments/$(tfplan)
	cd deployments && terraform apply "$(tfplan)" \
		&& rm -f "$(tfplan)"
