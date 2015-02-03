package main

var intro = `
STS: Secure Tunnel Server

Author: Xing Xing <mikespook@gmail.com>

Configuration file is using YAML format.

Here is an example:

addr: 127.0.0.1:2222
keys:
	id_rsa_1
	id_rsa_2
	id_rsa_3
log:
	file: /var/log/stsd.log
	level: all

And the meaning of fields are:

	* addr - What address and port stsd will bind to
	* keys - Private keys used by server
	* log  - Log configuration
	* log-file  - Log file, leave this empty to output to stdout
	* log-level - Log level: 'error', 'warning', 'message', 'debug', 'all' and
		'none',	levels can be combined with '|'

Running stsd:

$ stsd -config=config.yaml
`
