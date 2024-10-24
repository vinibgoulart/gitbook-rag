# gitbook-rag

This project is a proof of concept to show how to use the pgvector extension to make queries in a database using the cosine similarity.

## Brief explanation

The project have a sync service that makes a request to the gitbook api every 1 hour to get the data from the pages and save it in the database. The data is saved in the table `page` and the vector representation of the data is saved in the table `page.embedding`.

When you run the api or cli, you can make queries to the database using the cosine similarity to get the most similar pages to a given prompt.

![image](https://github.com/user-attachments/assets/71373c28-2fcd-4f32-bf81-b29b725dda0d)
![image](https://github.com/user-attachments/assets/474e8d4b-692f-410d-bccb-c1a230e733ae)

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

## Running with docker

```bash
docker-compose up -d
```

This command will start the postgresql and the pgadmin services.

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

## Running the api

To consume the database make queries to the vector index you need to run the api with the following command:

```bash
make api
```

### API Endpoints

| Method | Endpoint | Description                               | Parameters             |
| ------ | -------- | ----------------------------------------- | ---------------------- |
| POST   | /ai/page | Send a prompt to the chat bot             | { "prompt": "string" } |
| GET    | /chat    | Get the current session and chat messages | null                   |
| POST   | /logout  | Logout from the current session scope     | null                   |
