# go-course
Go-advance course on [Yandex.Practicum](https://practicum.yandex.ru/go-advanced/)

# How to

## Codegen and build binaries

```bash
make generate build
```

## Run container with DB and apply migrations
```bash
make up migrate
```

## Run server
```bash
make server
```

## Run agent
```bash
make agent
```

## Run tests
```bash
make test cover cover-html
```
Don't need DB for tests
Logged errors about crash worker it's expected behavior

## Run docs
```bash
make docs
```

## Run linter
```bash
make fix lint
```
