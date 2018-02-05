# aws-mfa-cli
Tiny Go program to set up a temporary AWS profile for MFA + CLI access

![Usage GIF](https://user-images.githubusercontent.com/379404/35804408-ae64c8f4-0a45-11e8-93f8-4f4c7ee7018a.gif)

## Installation

#### go get

```
$ go get -u github.com/Preetam/aws-mfa-cli
```

#### Binary

Binaries are available on the [Releases](https://github.com/Preetam/aws-mfa-cli/releases) page.

## Usage

```
$ aws-mfa-cli -h
Usage of aws-mfa-cli:
  -profile string
    	Profile name (default "default")
  -profile-mfa string
    	Temporary MFA profile name (default "profile-mfa")
  -region string
    	AWS Region (default "us-east-1")
  -token string
    	MFA token
```

## License

MIT
