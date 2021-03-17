# makefile

# file to speed up development and usage

# build the binary
build:
	$(clean_command)
	go build -o pre-commit

# package the binaries for different operating systems
package:
	$(clean_command)
	go test -v
	rm -rf releases/
	mkdir releases/
	GOOS=linux GOOARCH=amd64 go build -o pre-commit-x86-64-linux-1-0-1
	GOOS=darwin GOOARCH=amd64 go build -o pre-commit-x86-64-darwin-1-0-1
	GOOS=windows GOOARCH=amd64 go build -o pre-commit-x86-64-windows-1-0-1
	mv pre-commit* releases/

# build and run the binary
run:
	$(clean_command)
	go build -o pre-commit
	./pre-commit
	rm -f pre-commit

# run the test suite
test:
	$(clean_command)
	go test -v ./...

#EOF