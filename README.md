# gophkeeper
Client-server application to securely store and fetch user's secret data such as credential pairs, card information, any text data and any binary data
## Client
Client is TUI application that allows authentication, creating, editing secrets, file download and upload
## Server
Server is gRPC API using PostgreSQL as storage for user's data
# Run server
To build and run server docker containers:
`make up`
Check if all containers are running:
`make ps`
Stop containers:
`make down`