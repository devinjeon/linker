build: src
	cd src && \
		GOOS=linux go build && \
		zip ../deploy.zip linker && \
		rm linker
clean:
	rm -f src/linker deploy.zip
