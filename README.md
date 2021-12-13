# Licensebot Client for Go

## Update

Run `licensebot update [bundlePath]` to compile the repository's dependency
licenses into a JSON file at `bundlePath`. `bundlePath` defaults to
`licenses.json` in the repository's top level.

Licensebot will fail with an error if any dependencies are blacklisted.
This error is unavoidable by design.

This command is intended to be used during development, and the bundle file
should be committed alongside any dependency modifications.

## Check

Run `licensebot check [bundlePath]` to verify that the repository's license
bundle is up-to-date. `bundlePath` defaults to `licenses.json` in the
repository's top level.

This command is intended to be used in an automated build process, as a sanity
check to ensure the license bundle is not out of date.
