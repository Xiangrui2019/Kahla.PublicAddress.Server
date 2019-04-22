![KahlaLogo](https://raw.githubusercontent.com/AiursoftWeb/Kahla.App/dev/src/assets/144x144.png)

# Kahla.PublicAddress.Server

Kahla Public Address is Extend From Ganlvtech go-kahla-notify.

## How to Use

Before starting the server, you need to configure config.json.

After the configuration is completed, you can use the following command to start the public number server.

```bash
$ ./Server # Start public address server
```

## How to Config

To configure this server, you must configure config.json in the server startup directory

```json
{
  "PublicAddressName": "<Public account name>",
  "Email": "<Public number email>",
  "Password": "<Public number admin login Password>",
  "Port": <Public number api port>, 
  "CallbackURL": "<Public number callback server address>",
  "TokenStorageEndpoint": "<Public number Token Storage On Callback Server Endpoint>",
  "MessageCallbackEndpoint": "<Public number Message Callback On Callback Server Endpoint>"
}
```

## How to run

Before running, you need to install the dependencies:

```bash
$ ./installdeps.sh
```

You can run this project by running the following shell command:

```bash
$ go build *.go
$ ./Server
```

## Project Dependencies

[Golang](https://golang.org/)

We developed this under `Golang 1.12` and we recommend downloading the latest version of Golang 1.12 directly.

This project supports Windows Mac Linux.

## Project Dependencies SDK Download Address

If you are a normal user, you only need to install `installdeps.sh`.

## How to build exe

```bash
$ go build *.go
```

This command requires all the packages above.

The files under the ./ directory are the compiled binarys.

## How to build Linux Binrary

```bash
$ go build *.go
```

This command requires all the packages above.

The files under the ./ directory are the compiled binarys.

## Document

For more info please view [Kahla Wiki](https://wiki.aiursoft.com/ReadDoc/Kahla/What%20is%20Kahla.md).

## How to contribute

There are many ways to contribute to the project: logging bugs, submitting pull requests, reporting issues, and creating suggestions.

Even if you have push rights on the repository, you should create a personal fork and create feature branches there when you need them. This keeps the main repository clean and your personal workflow cruft out of sight.

We're also interested in your feedback for the future of this project. You can submit a suggestion or feature request through the issue tracker. To make this process more effective, we're asking that these include more information to help define them more clearly.