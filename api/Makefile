build: src/main.go
	cd src && \
		GOOS=linux go build -o linker main.go && \
		zip ../deploy.zip linker && \
		rm linker
clean:
	rm -f src/linker deploy.zip
