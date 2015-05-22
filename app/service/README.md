# Data Models Service

A service for programmatic consumption of the data models. The service is publicly hosted here: http://data-models.origins.link

## Build

**Dependencies**

- [Git](https://git-scm.com)
- [Go 1.3+](http://golang.org) ([test your installation](http://golang.org/doc/install#testing))

**Run**

```bash
make install build
```

This will put the `data-models` binary in your `$GOPATH/bin` directory. The examples below assumes `$GOPATH/bin` has been added to your `$PATH`.

## Usage

Once the binary is built, running the binary without any options will start the service. The service clones the data-models repository from GitHub into a `data-models` directory in your working directory. It is recommended to print the usage message to see the available options.

```bash
data-models -help
```

## Docker

Use the pre-built image on Docker Hub.

```bash
docker run -it -p 8123:8123 dbhi/data-models
```

Or build the image locally.

```bash
docker build -t dbhi/data-models .
```
