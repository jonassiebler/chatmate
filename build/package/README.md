# Package Configurations

This directory contains packaging configurations for different platforms and package managers.

## Structure

- `brew/` - Homebrew formula configurations
- `deb/` - Debian package configurations  
- `rpm/` - RPM package configurations
- `docker/` - Docker container configurations
- `snap/` - Snap package configurations

## Usage

Build scripts in `/scripts` will reference configurations in this directory when creating packages for distribution.
