# Delivery Service

## Overview

The delivery service sends a list of campaigns to a particular user meeting the targeting rules in terms of app, os and country.

## Requirements

Ensure the following are installed:

- **Go 1.23 or later**
- **Docker**
- **Docker Compose**
- **Mongo DB Atlas Cluster**  (MongoDB streams work only on replica set and Atlas provides without much hazzles).

## Setup

To set up and use this project, follow these steps:

1. **Clone the Repository:**

   ```bash
   git clone https://github.com/Ruthvik10/targeting-engine.git
   cd targeting-engine
   ```
2. **Follow the app.env.example and a create an app.env file your root directory with the actual values.**

2. **Build and Start Docker Containers:**
    ```bash
    make compose-up
    ```
3. **Stop and Remove Docker Containers:**
    ```bash
    make compose-down
    ```


## Makefile Targets

### `test`

Run the unit tests for the project, excluding specific directories:

```bash
make test
```

### `run`

Runs the service locally

```bash
make run
```

### `compose-up`

Starts the delivery service, mongo test container and redis as containers managed by docker compose.

```bash
make compose-up
```
### `compose-down`

Brings down the containers started by compose-up.

```bash
make compose-down
```
## Database scheme

The **schema** folder contains all the database related changes that needs to be applied to the database.

### `campaign_coll_schema.js`

Creates the campaigns collection with the json schema validation.

### `indexes.js`
Contains all the database indexes that needs to be applied.

### `seed.js`
Contains the seed data.

### `query.js`
Contains the query that can be used to retrieve the data from the database.
