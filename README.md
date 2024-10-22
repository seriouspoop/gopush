# Gopush

CLI tool for testing and pushing codes. This is aimed to be a free credential manager for git with some hint of CI tooling.

Manage your authentication token for `https` and key pairs for `ssh` all with gopush.
Gopush generated your `ssh` keys and give them to you for uploading to your remote provider i.e GitHub, BitBucket etc.

## Installation

Ensure that you have go installed, If not then visit [go installation](https://go.dev/doc/install).

```
go install github.com/seriouspoop/gopush@latest
```

Check gopush version to verify installation.

```
gopush -v
```
Add `$(go env GOPATH)` to the path.

## Usage

For non initialized git repo use

```
gopush init
```

For repos initialized with gopush use

```
gopush run
```

That's it, gopush will handle the rest
