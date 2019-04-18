# Backend
## Requirement
golang 1.11+

## How to Run
```sh
$ go build
$ ./main.go
```

## How to Run Postgrel in Docker
```sh
$ docker create -v /var/lib/postgresql/data --name PostgresData alpine
$ docker run -p 5432:5432 --name yourContainerName -e POSTGRES_PASSWORD=yourPassword -d --volumes-from PostgresData postgres
```