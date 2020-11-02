# makefile

# file to speed up development and usage

# create test files

create_run_files:
	$(clean_command)
	declare -a RUNTESTFILES=("./inventories/development/secrets.yml", "./inventories/development/default.yml")
	for i in "{RUNTESTFILES[@]}"; do
		touch RUNTESTFILES[$i];
	done

# package the binaries for different operating systems
package:
	$(clean_command)
	GOOS=linux GOOARCH=amd64 go build -o pre-commit-x86-64-linux-0-0-1
	GOOS=darwin GOOARCH=amd64 go build -o pre-commit-x86-64-darwin-0-0-1
	GOOS=windows GOOARCH=amd64 go build -o pre-commit-x86-64-windows-0-0-1

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