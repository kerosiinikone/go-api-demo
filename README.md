# Simple Go Demo API

My first simple Go demo API that uses ozzo-dbx ORM and Go HTTP server to construct a simple JSON API with JWT

Tests also included

Routes include:

```
GET /
GET /items (auth)
POST /items
```

To run with Docker:

```
make docker-run
```

To run tests:

```
make test
```

### TODO

Mock DB test
