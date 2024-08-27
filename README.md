# gophkeeper
Client-server application to securely store and fetch user's secret data such as credential pairs, card information, any text data and any binary data
## Client
Client is TUI application that allows authentication, creating, editing secrets, file download and upload
# Build client
To build client simply run
`make build`
# Run client
Depending on your OS, on linux:
`./bin/gophkeeper-linux`
on OSX:
`./bin/gophkeeper-darwin`
on Windows:
`./bin/gophkeeper-windows`
## Server
Server is gRPC API using PostgreSQL as storage for user's data
# Run server
To build and run server docker containers:
`make up`
Check if all containers are running:
`make ps`
Stop containers:
`make down`