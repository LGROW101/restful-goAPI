## Getting Started

```
https://github.com/LGROW101/restful-goAPI.git

cd restful-goAPI

docker compose up -d

Install swagger

go install github.com/swaggo/swag/cmd/swag@latest

http://localhost:8080/swagger/index.html

go run main.go

```

## RESTful API

# Create User

```
curl -X POST -H "Content-Type: application/json" -d '{"name":"John Doe",  "email":"john@gmail.com"}' http://localhost:8080/users

```

# GET USER

```
curl -X GET http://localhost:8080/users

```

# GET ID USER

```
curl -X GET http://localhost:8080/user/1

```

# PUT Updated USER

```
curl -X PUT -H "Content-Type: application/json" -d '{"name":"Updated John Doe","email":"updated_John@gmail.com"}' http://localhost:8080/users/id

```

Replace `id` with the actual user ID.

# DELETE id USER

```
curl -X DELETE http://localhost:8080/users/id

```

Replace `id` with the actual user ID.
