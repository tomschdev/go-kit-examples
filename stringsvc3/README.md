# stringsvc3

## Proxying

Diverting the execution of a service to another service. 
Proxying can be implemented as middleware, just like logging and instrumentatoin. \
The same decorator pattern is applied where a new concrete type, that wraps the previous concrete type, is introduced: 

```go
type proxymw struct {
	next StringService // serve most requests via this service
	uppercase endpoint.Endpoint // serve uppercase request via this endpoint
}
```


