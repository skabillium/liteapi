# LiteAPI example

This server implements the spec of the LiteAPI that was provided to me as part of an assessment from
Nuitee.

The purpose of this service is to receive requests like the following
`http://localhost:8080/hotels?/hotels/?checkin=2024-03-15&checkout=2024-03-16&currency=USD&guestNationality=US&hotelIds=129410,105360,106101,1762514,106045,1773908,105389,1790375,1735444,1780872,1717734,105406,105328,229436,105329,1753277&occupancies=[{"rooms":2, "adults": 2}, {"rooms":1, "adults": 1}]`

And after communicating with the HotelBeds API, return responses as specified.

## Setup
To setup the project you will need the following dependecies
- Make
- Go compiler

To install all dependencies run `make install`. After that you can build the binary with `make build`.
Please make sure to create a `.env` file with your specific environmet config following the spec from
the `.env.example` file.

## Project structure
The start of the server's lifetime is in `cmd/main.go` where the environment is loaded and validated
and the server is started. 

All functionality specific to the HotelBeds API is located at `cmd/clients/hotelbeds.go`. The purpose
of this structure is to separate the concerns between the usage and the internal business logic needed
to perform requests to the HotelBeds API.

In the `cmd/api` directory you can find all the endpoint handers of the server. These include:
- Documentation
- Health checks
- The `GET /hotels` endpoint that was required for the assignment


## Testing
For testing the server you can use the Postman collection provided in the assessment.
