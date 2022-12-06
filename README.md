# Player Update Tool

## About the Tool

This tool will call each client in order to eventually update its applications.
When running in a CLI, it takes 2 arguments:

- [1] the file path to the client ids (.csv)
- [2] a secret for the JWT

## Getting started

Clone the following repo:
`[https://github.com/sylvain-gdk/player-tech-assignment.git]`

### Start the update process

```bash
go run main.go <file path> <secret>
```

### Tests

```bash
go test -v
```
