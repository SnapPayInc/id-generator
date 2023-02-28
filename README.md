<img src="./LogoVinDrLab.png" width="256"/>

# VinDr Lab / IDGen

VinDr Lab IDGen (or ID Generator) is a part of the VinDr Lab project. Like the name it is, this service generates a new integer value which based on the key as string. 

## What does this project do?

The ID Generator service works like a counter service base on specified keys. Currently, it allows to increase by one and reset a counter key.

## Project tree

```
.
├── conf/ // configuration files
│   ├── config.development.toml
│   └── config.production.toml
├── constants/ // some constants for project
│   ├── api.go
│   └── errors.go
├── Dockerfile
├── generator/ //main processor
│   ├── generator_api.go
│   └── generator_db.go
├── go.mod
├── go.sum
├── LICENSE.md
├── main.go //runner
├── README.md
└── utils/ // utilities for project
    └── log_utils.go
```

## Installation

**Option 1: Kubernetes**

Go to deployment project and follow the instruction

**Option 2: Docker**

You can execute the <code>docker-compose.yml</code> file as follow:

```
docker-compose pull
docker-compose down
docker-compose up -d --remove-orphans
```

**Option 3: Bare handed**

Just run the main.go file

```bash
export GO111MODULE=on
go mod tidy
go run main.go
```

## Configuration

Following the Installation, the application has two ways to read the configurations. Once is from the <code>config.produciton.toml</code> file that comes with the app. Or you can override it by passing through environment variables in Docker.
As you can see, the configuration file has the following form:

```
[rqlite]
uri = "YOUR_RQLITE_URI"
```

Please note that, the conversion from environmental variables to API configuration items itself like: <code>RQLITE\_\_URI</code> equals to <code>rqlite.uri</code>

## Testing

Make some basic calls key as 'test' first.

```
http://localhost:8080/id_generator/test/tap (PUT)
```

Normally, it will return 1

```json
{"last_insert_id":1}
```

If you keep make request, it return:

```json
{"last_insert_id":2}
```

If you want to make a batch request, like this:

```
http://localhost:8080/id_generator/test/tap?count=5 (PUT)
```
it return:
```json
{"last_insert_ids":[7,6,5,4,3]}
```

After run the reset request, we got:

```
http://localhost:8080/id_generator/test/set\ # (POST)
--header 'Content-Type: application/json' \
--data-raw '{"value":100}'
```
```
http://localhost:8080/id_generator/test/tap
```

It returns to the value you set

```json
{"last_insert_id":101}
```

## Others

**More information**

For a fully documented explanation, please visit the official document.

**Roadmap**

As I mentioned above, the IDGen service itself does the job of generating value, simple. We hope this can inherit some awsome features from AtomicInteger of Java. Welcome to join us./
