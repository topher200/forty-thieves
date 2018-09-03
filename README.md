## Installation

1. Install PostgreSQL 9.4.x

2. Install Go 1.4.x, git, setup `$GOPATH`, and `PATH=$PATH:$GOPATH/bin`

3. Create PostgreSQL database.
    ```
    cd $GOPATH/src/github.com/topher200/forty-thieves
    go get github.com/rnubel/pgmgr
    pgmgr db create
    ```

4. Run the PostgreSQL migration.
    ```
    pgmgr db migrate
    ```

5. Run the server
    ```
    cd $GOPATH/src/github.com/topher200/forty-thieves
    go run main.go
    ```


## Environment Variables for Configuration

* **HTTP_ADDR:** The host and port. Default: `":8888"`

* **HTTP_CERT_FILE:** Path to cert file. Default: `""`

* **HTTP_KEY_FILE:** Path to key file. Default: `""`

* **HTTP_DRAIN_INTERVAL:** How long application will wait to drain old requests
    before restarting. Default: `"1s"`

* **DSN:** RDBMS database path. Default:
    `postgres://$(whoami)@localhost:5432/forty-thieves?sslmode=disable`

* **COOKIE_SECRET:** Cookie secret for session. Default: Auto generated.


## Running Migrations

Migration is handled by a separate project:
[github.com/rnubel/pgmgr](https://github.com/rnubel/pgmgr).

Here's a quick tutorial on how to use it. For more details, read the tutorial
[here](https://github.com/rnubel/pgmgr#usage).

```
# Installing the library
go get github.com/rnubel/pgmgr

# Create a new migration file
pgmgr migration {filename}

# Migrate all the way up
pgmgr db migrate

# Reset to the latest dump
pgmgr db drop
pgmgr db create
pgmgr db load

# Roll back the most recently applied migration, then run it again.
pgmgr db rollback
pgmgr db migrate

# Show the latest migration version
pgmgr db version
```


## Vendoring Dependencies

Vendoring is handled by a separate project:
[github.com/tools/godep](https://github.com/tools/godep).

Here's a quick tutorial on how to use it. For more details, read the readme
[here](https://github.com/tools/godep#godep).

```
# Save all your dependencies after running go get ./...
godep save ./...

# Building with godep
godep go build

# Running tests with godep
godep go test ./...
```
