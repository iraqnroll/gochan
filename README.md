# Gochan - A modern imageboard engine written in Go.

## About

Gochan is a free, lightweight, (hopefully) fast and (HOPEFULLY) user-friendly imageboard engine inspired by projects such as [Tinyboard](https://github.com/savetheinternet/Tinyboard) and [jschan](https://gitgud.io/fatchan/jschan).

## Requirements
> **Note:** Gochan is in very early development. Non-Docker build instructions are not provided at this stage.

To run Gochan, you'll need:

- Docker (version 20.10 or higher recommended)
- Docker Compose (version 2.x)

Optional for accessing from other machines or hosting on development a server:

- Open relevant ports for Gochan **(default 3000)** and Adminer **(default 3333)** in your firewall.
- Map your server's host address in your client machine to a temporary hostname **gochan.local**, this is needed for CSRF trusted origins.

## Usage
By default, Gochan ships with a preconfigured admin account:

- **Username:** `admin`
- **Password:** `admin`

For security reasons, you should change these credentials immediately after your first login - especially if you plan to host Gochan on a publicly accessible server.

## Support

>*Gochan is still in a very early stage of development — expect breaking changes and bugs.*

If you encounter bugs or broken functionality, please report them by creating an issue.  
Suggestions are welcome, but bug fixes will be prioritized until the project is more mature.