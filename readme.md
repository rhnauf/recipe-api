# Backend Assessment

---
<h3>Prerequisites:
- [Go](https://go.dev/dl/) programming language
- [PostgreSQL](https://www.postgresql.org/)
- makefile (optional)

---
I have attached the ddl script under ```./external/db/ddl.sql```
you can run it using migration tools like [goose](https://github.com/pressly/goose) or run it directly on the sql shell
---
copy the env example and modify to your own environment

```cp .env.example .env```

<br>
to install the dependencies

```go mod tidy```

<br>
to run the application

```make run``` or ```go run ./cmd/app```

<br>
alternatively if you prefer to build binary and run

```make build-run``` or ```go build -C ./cmd/app/ -o ../../recipe-api.exe && ./recipe-api.exe```

<br>
to run the test 

```make test``` or ```go test ./...```

---
for this assessment, I create basic CRUD recipe api, it is pretty minimum but I think it is enough to cover the basic requirements.
There are couples of point that we can extend based on the link provided, such as:
1. adding image & video for each recipe,for this we can use blob storage like AWS S3 to store the data and refer the metadata on the database
2. separate the ingredient, testimony & instruction tab to each their table, we can normalize by separating it to each table instead of storing long text on the ```recipes``` table
3. create rating functionality, for this we can create new table to store users rating regarding the recipe then calculate it based on the average of users rating