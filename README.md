# gitbook-llm

## Pre-requisites

- create the database in postgresql
- run the postgresql
- fill the environment variables in the .env file
- have the pgvector extension installed

## Installation

```bash
go mod tidy
```

## Running the server

To sync the data from the database to the vector index you need to run the server with the following command:

```bash
make server
```

With this, a request will be made to the gitbook every 1 hour to sync the data.

## Running the cli

To consume the database make queries to the vector index you need to run the cli with the following command:

```bash
make cli
```
