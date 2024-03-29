# CatFacts

Serverless Go project that utilizes Twilio to send a random cat fact to a list of subscribers once a day.

## Architecture

The project is structured as a single Go module with several packages. Executables (AWS Lambda handlers) are placed in a named directory under `cmd`. The root of the project is a package (`catfacts`). Domain types and functions are defined in this root package. `facts`, `sms` and `subscriptions` make up the domain types. Each additional package represents an integration between the domain types, and a particular service.

-  `catfactninja` provides a client used to fetch facts from an API.
-  `dynamodb` provides DynamoDB storage capabilities.
-  `lambda` provides Lambda handlers, helpers and middleware.
-  `twilio` provides a client to send SMS with twilio, as well as the capability to marshal strings into Twiml.

Structuring the application this way allows our types to semantically represent themselves. It becomes obvious which types belong to which domain (e.g. `catfacts.Subscription` or `dynamodb.SubscriptionStore`). Each package builds on the root package and acts as a thin integration layer between our core domain and a particular service or concept.

### Inspiration

This structure was heavily inspired by Ben Johnson's articles. Particularly [Packages as Layers](https://www.gobeyond.dev/packages-as-layers/).

## Structure

 ```
 .
├── Makefile                # Build script
├── go.mod
├── go.sum
│
├── cmd                     # Executables.
│   └── lambda
│       ├── send_fact       
│       │   └── main.go     # send_fact Lambda entrypoint.
│       └── twilio_callback
│           └── main.go     # twilio_callback lambda entrypoint.
│
├── catfactninja            # CatFactNinja client package.
│   └── client.go           # CatFactNinja client.
│
├── dynamodb                # Dynamo integration package.
│   └── subscriptions.go    # Subscription store.
│
├── lambda                  # Lambda integration package. Helpers, handlers and middleware.
│   ├── apigateway.go       # APIGatewayHandler interface definition.
│   ├── cloudwatch.go       # CloudWatchEventHandler interface definition.
│   ├── errors.go           # Lambda error handling.
│   ├── fact.go             # Lambda handlers for SendFact.
│   ├── lambda.go           # Root Lambda logic.
│   ├── middleware.go       # Middleware definitions for handlers.
│   └── twilio.go           # Lambda handlers for TwilioCallback.
│
├── twilio                  # Twilio integration package.
│   ├── sms.go              # Twilio SMS sender.
│   └── twiml               # Twiml package.
│       └── marshal.go      # Functions for marshalling a string to Twiml.
│
├── facts.go                # Fact domain type and functions.
├── sms.go                  # SMS domain type and functions.
└── subscriptions.go        # Subscription domain type and functions.
 ```
 
## Functions
 
### twilio_callback

twilio_callback is an HTTP triggered Lambda function. The function URL is configured as a Twilio webhook and is sent a payload when a text message is sent to the CatFacts phone number. The function does rudimentary parsing of the request, and determines if the user is subscribing, unsubscribing, needs help or the request cannot be understood. The function then takes the appropriate action (e.g. deleting a users phone number from the subscription table in Dynamo) before sending a response. 
 
### send_fact
 
send_fact is a time triggered Lambda function that runs once a day. The function fetches a random fact from the catfact ninja api, loads a list of subscribers from DynamoDB and attempts to send each subscription the fact.
 
## Additional services
 
[CatFactNinja](https://catfact.ninja/)
