# go-service-template
Golang microservice template

#### Features
- flat directory structure:
  - app - application logic and http server
  - cmd - main executable
  - internal - internal libraries
  - model - domains and models
  - store - data store
- CLI based on `kingpin`
- internal logging based on `zerolog`
- http server:
  - based on `httprouter`
  - CORS enabled
  - timeouts enabled
  - context middleware
  - log_request middleware
- store based on `postgresql`
  - reconnect if database is not available
  - custom pool settings
  - simple protocol enabled

#### TODO
  - instrumenting, tracing
  - read/write store pools
  - repeat database transactions if case of transport error
  - long queries canel on shutdown
  - pass context with deadline to store repositories
