# RSS Aggregator

This is a simple RSS aggregator built with a focus on fetching, parsing, and displaying RSS feeds in a user-friendly format. It enables users to subscribe to multiple RSS feeds and view the latest content from all their favorite sources in one place.

## Features (planned)

- Fetch RSS feeds from multiple sources.
- Parse and display feed content such as titles, descriptions, and publication dates.
- Follow and unfollow RSS feeds that other users have added
- Fetch all of the latest posts from the RSS feeds they follow

## Setup using Docker

The easiest way to get this project up and running is via [Docker](https://www.docker.com/). See [docs](https://docs.docker.com/get-started/) to get started.

Next, clone the repo:

```bash
git clone https://github.com/kkosiba/rss_aggregator.git rss_aggregator
```

and create a `.env` file at the repo root with the content like this:

```dotenv
# Server
HTTP_SERVER_PORT=8000

# Database
POSTGRES_HOST=db
POSTGRES_PORT=5432
POSTGRES_DB=default
POSTGRES_USER=app_user
POSTGRES_PASSWORD=S0m3_s3cr37!
```

Once this is set up, run the following command (make sure `make` tool is available):

```bash
make run
```

It may take a while for the process to complete, as Docker needs to pull dependencies, build the app and run database migrations. Once all of this done, the application should be accessible at `localhost:8000`.
