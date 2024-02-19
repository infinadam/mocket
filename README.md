# Mocket: HTTP Mocking Server

Mocket is an HTTP mocking/proxy service built in golang.  It is intended to sit
in local and staging enviroments to provide a reliable way to elicit a response
from a third-party API, and to act as a proxy in cases where an unmocked
response is desired. 

## Features

Request mocking (duh).  Given a set of files in the scripts directory, Mocket
will respond to a matched HTTP request with pre-defined status codes, headers,
and response bodies.  Requests can even be scripted to time out, or to have
their underlying sockets dropped prematurely.

When mocking isn't enough, Mocket will also serve as a proxy for indicated
requests, allowing for a combination of testing patterns.

## Usage

Simply create a new request/response script in your desired directory.  These
scripts are standard JSON files (TBD: other formats?), and come in different
flavors: HTTP mocks, HTTP webhook triggers, HTTP passthroughs, and TCP scripts.
Each have their own syntax and behaviors that are described below.

### HTTP Mocking

The most basic feature of Mocket.  These types of scripts allow for an HTTP
request matching a specific profile to receive a "canned" response from the
service.

The basic format of an HTTP mocking JSON script is as follows:

```
{
    "request": {
        "verb": "get",
        "url": "/my/third/party/path?id=1234",
        "headers": {
            "content-type": "application/json",
            "content-length": "100",
        },
        "body": {
            "data": "This is my data.  It is not 100 characters long.  Sorry, W3C!"
        }
    },
    "response": {
        "status": "421",
        "headers": {
            "content-type": "text/html"
        },
        "body": "Your integration is in another castle."
    }
}
```

#### Regular Expression

Mocking scripts support regular expression in the `request` object so that a
single script can service multiple requests.  For instance, in the above
example, should you want `421` responses issued for _any_ `id` queried,
regardless of value:

```
"url": "/my/third/party/path?id={{/\\d+/}}"
```

You can provide simple regular expression and flags to the `{{...}}`
delimiters.  Capture groups are globally-defined for each script, allowing you
to reference them later:

```
{
    "request": {
        ...
        "url": "/users/{{/(?<user_id>\\d+)/}}"
        ...
    },
    "response": {
        ...
        "body": { "userId": "{{$user_id}}" }
        ...
    }
}
```

These can also be referenced by index:

```
{
    "request": {
        ...
        "url": "/users?type={{/(basic|paid)/}}"
        ...
    },
    "response": {
        ...
        "body": { "userId": 123, "type": "{{$1}}" }
        ...
    }
}
```

### HTTP Webhook Triggers

Coming Soon!

### HTTP Passthroughs

Coming Soon!

### TCP Scripting

Coming Soon!
