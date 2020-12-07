build: src
	cd src && \
		GOOS=linux go build && \
		zip ../deploy.zip linker && \
		rm linker
clean:
	rm -f src/linker deploy.zip
plan: build
	terraform plan -out=terraform.tfplan
apply:
	terraform apply "terraform.tfplan"
	rm -f deploy.zip
