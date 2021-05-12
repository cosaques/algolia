<a href='https://github.com/jpoles1/gopherbadger' target='_blank'>![gopherbadger-tag-do-not-edit](https://img.shields.io/badge/Go%20Coverage-95%25-brightgreen.svg?longCache=true&style=flat)</a>

# [Algolia] Technical test

Small application that parses and processes a large amount of data exposing the following endpoints through a REST API:

* `GET localhost:<port>/1/queries/count/<DATE_PREFIX>`: returns a JSON object specifying the number of distinct queries that have been done during a specific time range
* `GET localhost:<port>/1/queries/popular/<DATE_PREFIX>?size=<SIZE>`: returns a JSON object listing the top `<SIZE>` popular queries that have been done during a specific time range

## Motivation

This is a way of trying to index the log file ([sample file](https://www.dropbox.com/s/duv704waqjp3tu1/hn_logs.tsv.gz?dl=0)) so that we can easily get the distinct and most popular queries.

### Index structure

Taking in account that `<DATE_PREFIX>` (from the REST API) has a  **limited** number of possible values, the idea was to try to pre-handle (to index) all its possible values. That way the complexity of any request is O(1).

> **Note** that it would not be possible if `<DATE_PREFIX>` would be a time range with an arbitrary end/begin dates. In that case we might split the time in a binary tree (descending to a minimum allowed time precision) - that way any time range could be represented as a combination of this binary tree's nodes.

### Index storage

For that example all the indexes are stored in a memory to be the most performant in terms of requests to the provided [sample file](https://www.dropbox.com/s/duv704waqjp3tu1/hn_logs.tsv.gz?dl=0).

> **Note** that in case of bigger data sets we might use the file system to store the indexes.

## How to install

This assumes that you have Go installed and setup.

Run the following commands :

```bash
$ git clone https://github.com/cosaques/algolia.git
$ cd algolia/site
$ go run . -help
$ go run . -addr=":<PORT>" -file='<PATH_TO_TSV_FILE>'
```

You should see a log telling that server is running, e.g. :
```bash
$ go run . -addr=":5000" -file='/gists/hn_logs.tsv'
2021/05/12 01:51:01 Starting the webserver on  :5000
```

## API

Once everything is working fine you can see a dashboard on :

`localhost:<port>`

![Index Dashboard](https://github.com/cosaques/algolia/blob/main/site/img/dashboard.png)

It shows you in real time the number of indexed queries from a logs file. When this counter stops incrementing - it means that the logs file is completely indexed.

Normally the indexation should take several minutes.

> **Note** that you can start requesting data (by endpoints described at the begining) without waiting the end of indexation. You'll get the results based on already indexed queries.

## Scalability

One of the possible ways to scale this application is to have a facade that handles the initial logs file and split it into several ones each having non-intersecting queries. Thus the search engine will be spreaded among N servers. Each server will receive only the queries that other servers don't possess.

Taking in account that queries don't intersect among the servers:

* to get a distinct count of queries - the facade should request each server and then sum their responses;
* to get N most popular queries - the facade should request each server for N most popular queries and then choose the N most popular among them.

It could be achieved by taking a hash-sum of each query and by retainng the remainder (mod function) from its division by the number of servers. The mod function will give us a server this query should go to.

