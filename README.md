## Simple SSE Chat Application
This is a simple chat application that uses Server-Sent Events (SSE) and HTMX to create a chat application. Allows multiroom chat. So far, the app has only two rooms: `Friends` and `Family` and the messages are not persisted. 

This app is built using the following technologies:
- [Go](https://golang.org/)
- [Leapkit](https://leapkit.dev/)
- [HTMX](https://htmx.org/)
- [Server-Sent Events](https://developer.mozilla.org/en-US/docs/Web/API/Server-sent_events)



## Features
- Multiroom chat
- Real-time chat
- Simple and easy to use

## Video Demo

https://github.com/user-attachments/assets/ddf412c1-dad3-4705-a7f6-257e791ad045


## Running the app
To run the app, you need to have Go installed on your machine. You can download Go from [here](https://golang.org/dl/).

After installing Go, you need to install kit, which is the CLI tool for Leapkit. You can install kit by running the following command:

```bash
go install github.com/leapkit/leapkit/kit@latest
```

After installing kit, you have to install the dependencies for the app and set up. You can do this by running the following command:

```bash
go mod download
```

```bash
go run ./cmd/setup
```

After setting up the app, you can run the app by running the following command:

```bash
kit s
```

This will start the app on `http://localhost:3000`.


## Contributing
Contributions are welcome. Feel free to open an issue or a pull request.




