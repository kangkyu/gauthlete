# Gauthlete

**Under development**

Go library for using [Authlete](https://www.authlete.com/), features are limited
but simple to use.

https://so.authlete.com/getting_started?locale=en

### Environment variable
- `AUTHLETE_SERVICE_APIKEY`: "API Key" in "Service Details"
- `AUTHLETE_SERVICE_APISECRET`: "API Secret"

### How to develop
Go to authorization server code and run the go file with two environment variables below, as needed by Gauthlete library. And generate a client id and client secret from the authorization server. And then go to client app and run the go code by `go run .` as well as backend app.

Example apps (for development):
- client app https://github.com/kangkyu/gauthlete-test-client-app
- authorization server https://github.com/kangkyu/gauthlete-test-application
- resource server https://github.com/kangkyu/gauthlete-test-backend-app
