# stringsvc2
stringsvc2 separates stringsvc1 into respective files for better modularisation - this allows a go-kit project to maintain readability and separation of concerns as it scales. \
Therefore, we introduce: \
`service.go` for holding Service interface, service struct (concrete type), service methods to implement Service, error types. \
`transport.go` for holding endpoint adapter functions (one for each method in Service interface), request and response structs (one pair for each method in Service interface), decoding helper functions (one for each request struct), and encoding helper functions (if possible, one for all response structs). \
`logging.go` for decorating Service by wrapping StringService in a new concrete type that also holds a logger, and reimplementing Service interface with the new concrete type. The new implementations execute logging calls and then diverts execution to the original implementations of the interface (pure business logic), by using the instance of the original concrete type. \
`instrumenting.go` serves the same purpose as logging.go, except it adds prometheus observability as opposed to logging. \
`main.go` for stitching components together in func main(). This involves initialising concrete types and dependencies (db connection, logger etc), chaining middleware decorators on our service implementation, creating handlers by linking endpoints (via endpoint adapters) to transport, creating a router by mapping handlers to routes (in http case), and serving the server. \
Note: no global scope! for example, logger is attached to service via middleware.

## Middleware
Middleware layers functionality on top of existing functionality by decorating a service (wrap concrete type with new concrete type and reimplement methods, by referencing original implementations). \
This process is implemented by wrapping our concrete type in a new concrete type, and then reimplementing the original interface by first adding additional logic (the purpose of the middleware e.g., logging), and then invoking the implementation of the original concrete type by accessing it through the instance of the original concrete type held as a field in our new concrete type.

Middleware returns Endpoint returns Service invokes service. \
Service is resposible for executing business logic on inputs. \
Endpoint is responsible for extracting fields from requests, feeding to service and returning response and error. An endpoint looks after the provisioning of a service's method. \
Middleware is responsible for decorating endpoint with supplementary services by running additoinal logic and then diverting execution to endpoint. Additional logic is typically provided by a dependency i.e., logger, prometheus, authentication, etc. 

## Execution
```bash
cd stringsvc2
go build 
./stringsvc2
```
From new terminal session: \
```bash
curl -XPOST -d'{"s":"hello, world"}' localhost:8080/uppercase
curl -XPOST -d'{"s":"hello, world"}' localhost:8080/count
```
Notice: the terminal running the server will have outputted logs and instrumentation reports after processing the requests.
