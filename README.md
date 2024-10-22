# gitbook-llm

This project is a proof of concept to show how to use the pgvector extension to make queries in a database using the cosine similarity.

## Pre-requisites

- create the database in postgresql
- run the postgresql
- fill the environment variables in the .env file
- have the pgvector extension installed

### How to install the pgvector extension

```bash
git clone https://github.com/pgvector/pgvector.git
cd pgvector
make
sudo make install
```

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
