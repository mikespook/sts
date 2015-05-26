package main

const intro = `
STS: Secure Tunnel Server - MIT Software License 

	Xing Xing <mikespook@gmail.com>
	http://github.com/mikespook/sts

CONFIGURATION

	Config file should be in YAML format.

	Here is an example:

--------- Start ----------


log:
    file:
    level: all
pwd: ../../misc
rpc:
    addr: 127.0.0.1:2223
tunnel:
    addr: 127.0.0.1:2222
    auth:
        anonymous: false
        password: static://123456
        # password: rpc://127.0.0.1:9000
        pubkey: file://id_rsa.pub
        # pubkey: rpc://127.0.0.1:9000
    keys: [ id_rsa ]

---------- End -----------

	* tunnel-addr: "stsd" will bind to this address (and port)
	* tunnel-auth: See "AUTHORIZATION" for more details
	* tunnel-keys: Private keys used by server
	* log: Log configuration
	* pwd: Current working directory
	* rpc-addr: "stsd" will bind a RPC service to this address (and port)

AUTHORIZATION

	stsd has three auth-options: anonymous, password and pubkey.

	When "anonymous" is set to "I Understand the Risks"(DANGER!!!), other 
	options will be ignored and there will be no authenticated mechanism to stsd.

	"password" and "pubkey" should be one of following formats:

	* "static://[string]", the string will be used as password/pubkey directly;
	* "file://[file]", the file will be read and its contents will be used as 
		password/pubkey;
	* "rpc://[url]", connect to a rpc server witch can serve two remote calls:
		1) STS.PasswordAuth 2) STS.PublicKeyAuth. And "url" should be in form 
		"(tcp|http)://[ip]:[port]/[path]".

LOG

	In the default way, the output of log will print to stdout. In case of
	"file" is a legal path, Log will output records to it. Levels can be 
	combined with "|" from following values: "error", "warning", "message",
	"debug". Also, there are two magic values "all" and "none" to display all
	levels or nothing.

USAGE

	Running stsd:

	$ stsd -config=config.yaml

	Establish local ssh tunnel proxy:

	$ ssh -p 2222 127.0.0.1 -D18081

	Using SOCKS5 proxy:

	$ curl --socks5 127.0.0.1:18081 http://www.mikespook.com
`
