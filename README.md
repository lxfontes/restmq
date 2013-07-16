See [tunacorp's restmq](https://github.com/gleicon/restmq) page.


This is an drop-in replacement for the current cyclone implementation.

# Installation

    make

This is equivalent to:

    go get
    go build

Then run:

    ./restmq

# Supported operations

## /q/\<queue\>

*Generic handler*


- GET /q/\<queue\>

Take element from queue

- POST /q/\<queue\> with "value=zzz"

Add element "zzz" to queue

- DELETE /q/\<queue\>

Flush queue (delete it)


## /queue
*JSON handler*

Single endpoint for any operation. JSON commands are to be passed via 'body' parameter. See `tests` directory for examples.

- GET

Request

	{
		"cmd": "get",
		"queue": "test"
	}

Response
	
	{"queue":"test","value":"pizza","key":"test:1"}

Retrieve top entry from queue without deleting it

- DELETE

Request

	{
		"cmd": "del",
		"queue": "test",
		"key": "test:1"
	}

Response

	{"queue":"test","key":"test:1","deleted":true}
	
Deletes specific entry from queue

- ADD

Request

	{
		"cmd": "add",
		"queue": "test",
		"value": "pizza"
	}

Response

	{"queue":"test","key":"test:2"}

Push entry to queue

- TAKE

Request

	{
		"cmd": "add",
		"queue": "test",
		"value": "pizza"
	}

Response

	{"queue":"test","value":"pizza","key":"test:5"}

Combination of GET + DELETE as a single operation
