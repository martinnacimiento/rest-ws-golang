# rest-ws-golang

A RESTful API and WebSocket service built with Go. This project serves as a learning exercise, where I implement various features such as authentication, CRUD operations for posts, and real-time communication using WebSockets.


## 🚀 Getting started
To get started, ensure you have Docker installed. Then, follow these steps to set up the database and run the application:

1. Navigate to the `database` directory: 
```sh
cd database
```
2. Build the PostgreSQL Docker image:
```sh
docker build . -t postgres-rest-ws-golang
```
3. Run the PostgreSQL container:
```sh
docker run -p 54321:5432 postgres-rest-ws-golang
```
4. Return to the project root directory:
```sh
cd ..
```
5. Run the Go application:
```sh
go run main.go
```

Now, you can test the API endpoints using an HTTP client like Postman.