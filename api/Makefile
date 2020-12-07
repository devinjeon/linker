build: src
	cd src && GOOS=linux go build
clean:
	rm -f src/linker deploy.zip
plan: build
	terraform plan -out=terraform.tfplan
apply:
	terraform apply "terraform.tfplan"
	rm -f deploy.zip
