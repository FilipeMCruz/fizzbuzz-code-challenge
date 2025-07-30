# fizzbuzz-code-challenge

This document presents the needed information to install and run the solution given for the Code Challenge. To better
understand the code there are some comments in it.

This solution is a simple http server written in go.

Since the domain and requirements are simple there's no need to over-complexify the solution, therefore I choose to:

- not use any libraries, go already provides the basics needed to do what this challenge requires;
- not follow a commonly used architecture such as onion, layered or clean architecture for the sake of simplicity.

API details can be found [here](docs/API.md).
Architecture details can be found [here](docs/ARCHITCTURE.md)

## Build & Running

There's two different ways to run the solution:

- natively, requires golang v1.22+ to be installed;
- docker, requires docker to be installed.

Note that this has only been tested in linux.

### Natively

Ensure that the go compiler is available in your workspace.

To build the solution with go:

```shell
go build -o server .
```

To run the solution in port 8080:

```shell
./server -port 8080
```

### Docker

To run the solution in port 8080:

```shell
docker compose up -d
```

## Requirements assumed

- HTTP responses are provided in json since it is probably the most used format for modern HTTP services;
- HTTP requests for the fizzbuzz endpoint define all required params in the query string;
- the stats endpoint only cares about the defined endpoints (/stats and /fizzbuzz);
- the stats endpoint is simple and assumes that query param order matters (/example?a=1&b2 is different from
  /example?b=2&a=1).

## Possible improvements

This solution assumes that a single instance running is enough and there's no need to keep the stats for ever.
If that's false, I'd suggest to horizontally scale the solution and check if a cache would help tackle the performance
issues.

Horizontal scaling provides flexibility when handling different volumes of data but does require the following changes:

- a load balancer needs to be placed in front of the solution (I'd probably try to use the 'Least Connection' algo since
  fizzbuzz is purely a CPU intensive task);
- the stats have to be shared and stored in disk (I'd probably use redis to keep track of all requests made).

## Notes

I don't know how 'production ready' this would need to be since there's no information about:

- where it would run (single container, serverless function in the cloud, on-prem in a VM...);
- expected requests per second;
- authentication/authorization needs;
- accessibility of the service, e.g. will it be publicly exposed? If so, how are TLS/SSL certificates normally used
  within the company?
- infrastructure tied to it (for observability/monitoring, distributed logging, orchestration, etc..);
- commonly used libraries within the company (e.g. gin, testify);
- what tools are used to document the api surface (e.g. swagger/OpenAPI, simple API.md);
- what linting rules are used;
- standard encoding and data formats used when exchanging data between http services within the company;
- how API versioning is tackled.

I welcome any discussion on these topics and others that I may have forgotten about.
