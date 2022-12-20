# PasswordServer2

## Design
### Frontend
Frontend is located in `/frontend`. The frontend is built with SvelteKit and is embed within the built go executable.

### Backend
Backend is go using the standard `net/http` library. 

## Usage
### Running
`go generate ./... && go run server.go`

### Building
```
go generate ./... # Build frontend
go build server.go # Build backend

./server.go # Run resulting executable
```