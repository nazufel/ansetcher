# Ansible Secrets Watcher

Ansible Secrets Watcher is a utility to run as a [git pre-commit hook](https://git-scm.com/docs/githooks#_pre_commit) to check for plaintext secrets files in an [Ansible](https://docs.ansible.com/) repository that should be encrypted with [Ansible Vault](https://docs.ansible.com/ansible/latest/cli/ansible-vault.html) before running a git commit in order to protect secrets from being committed into a source repo in plaintext.

## Install

Users can either pull down a prebuilt binary for their operating system or build from source.

### Using a Prebuilt Binary
 
The [MAINTAINERS](./MAINTAINERS.md) have prebuilt operating-system specific binaries and provided them in the [releases](./releases) directory for Linux, MacOS, and Windows on the amd64 CPU architecture. The binary files are named as such for each operating system, architecture, and release. Other operating systems or CPU architectures can be either compiled by the user that needs it, or issues opened in this repository requesting for additional architectures or operating systems are welcome.

If a user decides to pull down a prebuilt binary, they will still need to copy the file to their own ```.git/hooks/``` directory for each repository that has Ansible secrets.

### Compiling from Source

Ansible Secrets Watcher is written in [Golang](https://golang.org) (Go) and tested against version ```1.15.3```. To build from source, make sure you have [Go](https://golang.org/dl/) installed at the specified version or higher. If you have [GNU Make](https://www.gnu.org/software/make/) installed, then you can use the included [makefile](./makefile) to build the binary for you with the command, ```make package```. This will run all of the tests. create a ```releases/``` directory, and build three binaries for three different operating systems. This is the same process used to create the binaries above. 

If you do want to run the build command yourself for reasons including, but not limited to, you need a different operating system or CPU architecture than what is provided, then you can run the following command:

```sh
go build -o pre-commit
```

This will build all of the ```*.go``` files into a single binary named ```pre-commit```. The user will then still need to copy it to their ```.git/hooks/``` directory in each repo they wish use this utility. 

By default, the ```go build``` command will use the calling system's operating system and CPU architecture to compile the binary to the that system. If the user wishes to compile to a different operating system or architecture, then the use must look up the desired target system in the Go [syslist](https://github.com/golang/go/blob/master/src/go/build/syslist.go) to see if it is supported. If so, they will just need to set to environment variables in the ```go build``` steps: ```GOOS``` and ```GOARCH```. For example, to compile for FreeBSD on a RiscV 64-bit system would look like this:

```sh
GOOS=freebsd GOARCH=riscv64 go build -o pre-commit
```

## Using

The Ansible Secrets Watcher utility runs a pre-commit hook, as discussed above. Git will run any scripts or programs in the ```.git/hooks``` directory of a git repository. These will get called before the actual ```git commit``` is passed. If there's an error from the ```.git/hooks/pre-commit```, then ```git commit``` is aborted. This is the point. It allows to check for certain conditions before the commit is processed.

```sh
git commit -am "commting some awesome new functionality"
2020/11/02 14:26:47 found inventory root to walk: ./inventories, proceeding
2020/11/02 14:26:47 Error! Found Ansible Vault secrets file in plaintext during commit: inventories/development/secrets.yml. Please encrypt the file and reattempt to commit.
2020/11/02 14:26:47 Error! Found Ansible Vault secrets file in plaintext during commit: inventories/production/secrets.yml. Please encrypt the file and reattempt to commit.
```

The program existed with an error, which stopped the commit. Checking the status shows files with uncommitted changes:

```sh
git status
On branch main
Your branch is up to date with 'origin/main'.

Changes not staged for commit:
  (use "git add <file>..." to update what will be committed)
  (use "git restore <file>..." to discard changes in working directory)
        modified:   watcher.go
        modified:   watcher_test.go

no changes added to commit (use "git add" and/or "git commit -a")
```

There are still files showing they have changes that need to be committed. After properly Ansible Vault encrypting the files we are finally able to commit the changes.

```sh
git commit -am "commiting some awesome new functionality and encrypted vault secrets"
2020/11/02 14:39:31 found inventory root to walk: ./inventories, proceeding
[main a34c94f] commiting some awesome new functionality and encrypted vault secrets
 2 files changed, 19 insertions(+), 11 deletions(-)
 rewrite README.md (88%)
rross@gloomwalker ~/go/src/github.com/nazufel/ansible-secrets-watcher (main) $ git status
On branch main
Your branch is ahead of 'origin/main' by 1 commit.
  (use "git push" to publish your local commits)

nothing added to commit but untracked files present (use "git add" to track)
```