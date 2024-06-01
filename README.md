# REST API for Qlik demonstration

Welcome to my repo for the Qlik audition project.

## Introduction

(TODO - Explain overall setup with Docker, Go, and Postgres)

## Installation

 - Set up your local environment with Docker and Git
 - Clone this repository into a project directory on your computer
 - Copy the `example.env` file and rename it to `.env`
   - Populate with the values needed (some defaults are provided)
 - Run: `docker compose up`
   - This will create two images: `app` and `postgres`
 - Query the API:
   - Open a web browser to `localhost:8080` (changing `8080` to what you set for `API_PORT`)
   - CURL the endpoint in your terminal `curl http://localhost:8080` (changing `8080` to what you set for `API_PORT`)
