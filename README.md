# STS

STS: Secure Tunnel Server - a configurable and powerful ssh tunnel server.

# Motivation

# Installation

	$ go get github.com/mikespook/sts

# Configuration

Config file should be in YAML format.

Here is an example:

```yaml
addr: 127.0.0.1:2222
auth:
	anonymous: false
	password: static://123456
	pubkey: file://id_rsa.pub
keys: [ id_rsa ]
log:
	file:
	level: all
pwd: ../misc
```

* addr: stsd will bind to this address (and port)
* auth: See below
* keys: Private keys used by server
* log: See below
* pwd: Current working directory

## Authorization

stsd has three auth options: anonymous, password and pubkey.

When `anonymous` is set to `I Understand the Risks`(DANGER!!!), other options will be ignored and there will be no authenticated mechanism to stsd. `password` and `pubkey` should be one of following formats:

* `static://[string]`, the string will be used as password/pubkey directly;
* `file://[file]`, the file will be read and its contents will be used as password/pubkey;
* `rpc://[url]`, connect to a rpc server witch can serve two remote calls: 1) STS.PasswordAuth 2) STS.PublicKeyAuth. And "url" should be in form `(tcp|http)://[ip]:[port]/[path]`.

## Log

In the default way, the output of log will print to stdout. In case of `file` is a legal path, Log will output records to the file. Levels can be combined with `|` from: `error`, `warning`, `message` and `debug`. Also, there are two magic values `all` and `none` to display all levels or nothing.

# Usage

Running stsd:

	$ stsd -config=../misc/config.yaml

Establish local ssh tunnel proxy:

	$ ssh -p 2222 127.0.0.1 -D18081

Using SOCKS5 proxy:

	$ curl --socks5 127.0.0.1:18081 http://www.mikespook.com

# TODO

 * Session status: connections, R/W bytes, start/established time
 * Server status: active clients, R/W bytes, start/established time

# Contributors

(_Alphabetic order_)
 
 * [Xing Xing][blog] &lt;<mikespook@gmail.com>&gt; [@Twitter][twitter]

# Open Source - MIT Software License

See [LICENSE][license].

[blog]: http://mikespook.com
[twitter]: http://twitter.com/mikespook
[license]: LICENSE
