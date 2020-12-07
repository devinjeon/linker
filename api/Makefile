tfplan=terraform.tfplan

build: src
	cd src && GOOS=linux go build
clean:
	rm -f src/linker deploy.zip
plan: build
	terraform plan -out="$(tfplan)"
apply: tf-apply clean

tf-plan:
	terraform plan -out="$(tfplan)"
tf-apply: $(tfplan)
	terraform apply "$(tfplan)" \
		&& rm -f "$(tfplan)"
