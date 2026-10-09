# BehindTheLines

A multiplayer top-down game written in Go. The repository has a TCP game server and an Ebitengine graphical client.

## Requirements

- Go (the modules declare Go 1.15)
- A desktop environment supported by [Ebitengine](https://ebitengine.org/) for the graphical client

The client uses Go modules. `go run` downloads the required modules the first time it is run, so network access to the Go module proxy may be needed.

## Run the game

Start the server in one terminal:

```sh
cd server
go run .
```

Start **three client instances**, each in a separate terminal:

```sh
cd frontend
go run .
```

The server listens on TCP port `8080`. Each client connects to `localhost:8080`. The game starts after all three clients connect; clients may appear to wait before then. Run the client from the `frontend` directory because it loads its images and map using paths relative to that directory.

## Controls

- `W`, `A`, `S`, `D`: move
- Left mouse button: shoot
- `Q`: quit

## Project layout

- `server/`: Go TCP server and game/player logic
- `frontend/`: Ebitengine game client, networking, and assets
