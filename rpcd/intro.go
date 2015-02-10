package main

const intro = `
RPC Daemon for Secure Tunnel Server - MIT Software License 

	Xing Xing <mikespook@gmail.com>
	http://github.com/mikespook/sts/rpcd

CONFIGURATION

	Config file should be in YAML format.

	Here is an example:

---------Example Start----------

addr: tcp://127.0.0.1:9000
log:
	file:
	level: all
mongo:
	addr: 127.0.0.1
	db: sts	
pwd: ../misc

----------Example End-----------

	* addr: "rpcd" will bind to this address (and port)
	* log: Log configuration
	* mongo: 
		* addr: The address and port that MongoDB is listening on
		* db: Database name
	* pwd: Current working directory

LOG

	In the default way, the output of log will print to stdout. In case of
	"file" is a legal path, Log will output records to it. Levels can be 
	combined with "|" from following values: "error", "warning", "message",
	"debug". Also, there are two magic values "all" and "none" to display all
	levels or nothing.

USAGE

	Running stsd:

	$ rpcd -config=config.yaml

`
