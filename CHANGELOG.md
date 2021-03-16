# Changelog
All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [1.0.1] - 2020-03-16
### Added
- Two configuration options to change the behavior of the secrets watcher.

## [1.0.0] - 2020-11-02
### Added
- Initial program and tests [@nazufel](https://github.com/nazufel).
- program walks down the passed inventory root directory looking files named "secrets.yml"
- program checks the found "secrets.yml" files and makes sure they are encrypted with ansible vault, if not, throw an error with path to the plaintext file