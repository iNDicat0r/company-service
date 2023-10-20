# Company uService

## Framework and Libraries Used
I have used the following framework and ORM:
- Gin to make routing
- Gorm

The reason I used them is to speed up the development of this assignment.

## Quick start
From the root of the repository:

```
docker compose up

make migrate
```
and then to run the app:
```
make run
```
The following commands are available to run:
```
make build
make lint
make run
make migrate
make unit
make coverage
```


## Design and assumptions
The overall design is illustrated here:![Alt text](/docs/design.png "Design")

1. The design decisions I took for implementing this service was to separate the application into layers. I could follow a quick approach which is basically to handle the business logic in a controller-like structure and have the business logic implemented there. However I decided to decouple each part of the application to demonstrate a "domain driven" design with each layer dealing with their responsibility.

2. The models are used across different layers, while this approach allows reusability it also introduces some dependency and coupling. In the context of this assignment I used them across services and repositories.

3. All the application contexts (starting from routes) are propagated down to the layers and to database operation.

## Improvements
The following are a list of improvements that can be done:
- Migrations can be done in a standard way through external CLI
- Seeding can be done using specific command and utility
- Caching
- Domain types so services are completely decoupled from other packages. I have demonstrated this in the company service by implementing local types. This can be extended to other layers too.

