# Websocket Chat

Websocket Chat is a Golang pet project I chose to write to learn Golang. There is no frontend at the moment. I used this project only for educational purposes in backend development. Thus, I didn't really want to write any frontend for this project. But there'll be for other pet projects, because I didn't want to spend much time on Websocket chat and spend more time on learning Golang + backend. My next projects will be more ***representative***

### Motivation
My motivation was to learn how to build anything in Golang. I wanted to learn how to create HTTP server, do some CRUD operations, native auth with JWT, proper logging, config set up, a bit of unit testing, etc. I got the idea of writing exactly Websocket server from Primeagen's New Language Roadmap: **https://www.youtube.com/watch?v=E1H9AFtwxfE**

## Installation

To run it locally, use:

```bash
git clone https://github.com/LeoDiKadyrov/newgolang-websocketchat.git
go mod tidy
cd cmd/websocket-chat
go run main.go
```

Note: You'd have to set up a Postgres storage to run it. 

## Usage

To access Swagger API Doc first run project locally
```bash
go run main.go
```
Then navigate in browser to:
```
http://localhost:8080/swagger/index.html#/
```
![image](https://github.com/LeoDiKadyrov/newgolang-websocketchat/assets/60335678/c6b0e63e-13fc-46fa-bcfa-eca8149427f6)


## Structure
```
Folder Structure
/cmd/websocket-chat
│ └── main.go
/config
│ └── local.yaml
/internal
│ ├── /config
│ │ └── config.go
│ ├── /http_server
│ │ ├── /handlers
│ │ │ ├── /jwt - handlers for refreshing JWT tokens
│ │ │ │ └── /mocks
│ │ │ ├── /user - handlers for saving and deleting user from database 
│ │ │ │ └── /mocks
│ │ └── /middleware - custom middleware for slogger
│ │   └── /logger
│ ├── /lib
│ │ ├── /api - custom responses, errors
│ │ ├── /encryption - user password encryption (bcrypt)
│ │ ├── /jwt - generating, extracting from requests, validating JWTs
│ │ │ └── /middleware - custom middleware for authentication
│ │ ├── /logger
│ │ │ ├── /handlers
│ │ │ │ └── /slogdiscard - to remove logs during tests
│ │ │ ├── /sl - custom error func for slogging
│ └── /storage - storage set up
│ │ └── /postgres - postgres set up
│ ├── /websocket 
│ │ └── /handlers - upgrader and message handlers
/go.mod
/go.sum
/local.env
```
## Built With

* [Chi](https://go-chi.io/) - The router used
* [Gorilla Websocket](https://github.com/gorilla/websocket) - For setting up Websocket server
* [Playground Validator](https://github.com/go-playground/validator) - Used for data validation in HTTP requests
* [Clean Env](https://github.com/ilyakaznacheev/cleanenv) - To work with config
* [jwt-go](https://github.com/golang-jwt/jwt) - For using JWT as authentication method
* [Testify](https://github.com/stretchr/testify) - For testing
* [Mockery](https://github.com/vektra/mockery) - To generate mocks for testing HTTP handlers

## License

[MIT](https://choosealicense.com/licenses/mit/)
