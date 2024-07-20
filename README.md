# Project management

Project management REST API

## Features

- Creating tasks, users, projects
- Searching

## Installation & Usage

1. Navigate to the project directory: `cd project-management`.
2. Rename .env.example to .env and change variables accordingly.
3. Start the docker containers: `make up`.
4. Navigate to swagger docs at http://localhost:8080/swagger/index.html.

## Libraries

1. [go-chi](https://github.com/go-chi/chi) as router
2. [zerolog](https://github.com/rs/zerolog) as logger
