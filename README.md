# authentication
ClawIO Authentication Service

It contains 3 components:

| [authenticationcontroller](authenticationcontroller)  | [service](service)       | [server](server)       |
| ------------------------- | ------------- | -------------|
| [![Go Report Card](https://goreportcard.com/badge/github.com/clawio/authentication/authenticationcontroller)](https://goreportcard.com/report/github.com/clawio/authentication/authenticationcontroller)  | [![Go Report Card](https://goreportcard.com/badge/github.com/clawio/authentication/service)](https://goreportcard.com/report/github.com/clawio/authentication/authenticationcontroller) | [![Go Report Card](https://goreportcard.com/badge/github.com/clawio/authentication/authenticationcontroller)](https://goreportcard.com/report/github.com/clawio/authentication/server) |


This service exposes the following HTTP/2 endpoints:

## Authenticate 

### Request

```
GET http://localhost:58001/clawio/v1/auth/verify/<token>
```

### Response

HTTP Status Code: 200

```
{
	"username": "test",
	"email": "test@test.com",
	"display_name": "tester"
}
```

## Authenticate 

### Request

```
POST http://localhost:58001/clawio/v1/auth/authenticate
```
Body:

```
{
	"username": "test",
	"password": "testpwd"
}
```

### Response

HTTP Status Code: 200

```
{
	"token": "testoken"
}
```
