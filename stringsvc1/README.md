# stringsvc1

A basic introduction to go-kit's API building features, all centralised in `main.go`. \
Refer to comments in `main.go` for explanation of features.
The microservice implemented performs operations on strings, received in JSON request and decoded into request struct and inputted into one of the service's endpoints.

## Execution:
```bash
cd stringsvc1
go run main.go
```
From a new terminal session: \
Uppercase endpoint:
```bash
curl -XPOST -d'{"s":"hello, world"}' localhost:8080/uppercase
```
Count endpoint:
```bash
curl -XPOST -d'{"s":"hello, world"}' localhost:8080/count
```


