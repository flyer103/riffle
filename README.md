# Riffle - RSS Feed Reader and Aggregator

Riffle is a powerful RSS feed reader and aggregator with a REST API for managing RSS sources and content.

## Features

- **RSS Source Management**: Add, update, delete, and list RSS sources
- **Content Management**: Fetch, update, delete, and list RSS content
- **Recommendations**: Get personalized content recommendations based on user feedback
- **Search**: Search for content by keywords
- **Batch Operations**: Perform batch operations on sources and content
- **OPML Import**: Import RSS feeds from OPML files
- **Content Analysis**: Analyze RSS content quality and relevance
- **Metrics**: Prometheus metrics for monitoring
- **Profiling**: Optional pprof endpoints for debugging
- **Modern Web UI**: A responsive web interface built with Vue.js and Material 3 design

## Getting Started

### Prerequisites

- Go 1.24 or higher
- SQLite3
- Node.js 16+ and npm (for frontend)

### Installation

1. Clone the repository:

```bash
git clone https://github.com/flyer103/riffle.git
cd riffle
```

2. Build the backend application:

```bash
make build
```

3. Set up the frontend:

```bash
cd frontend
npm install
```

### Usage

Riffle provides several commands:

#### Running the Server

```bash
./riffle serve --port 8080 --db-path ./riffle.db
```

#### Running the Frontend

```bash
cd frontend
npm run serve
```

Then open your browser to http://localhost:3000

#### Importing OPML Files

```bash
./riffle import-opml --opml feeds.opml --db-path ./riffle.db
```

#### Analyzing RSS Feeds

```bash
./riffle run --opml feeds.opml --interests interests.txt --articles 5 --top 10
```

#### Command-line Options

##### Serve Command Options
- `--port`: Port to listen on (default: 8080)
- `--db-path`: Path to the SQLite database file (default: ./riffle.db)
- `--log-level`: Log level (debug, info, warn, error) (default: info)
- `--enable-pprof`: Enable pprof debugging endpoints (default: false)
- `--metrics-port`: Port for Prometheus metrics (0 to disable) (default: 0)
- `--rate-limit`: Rate limit in requests per second (0 to disable) (default: 100)
- `--enable-cors`: Enable CORS (default: false)
- `--cors-origins`: Allowed CORS origins (default: *)
- `--read-timeout`: HTTP server read timeout (default: 30s)
- `--write-timeout`: HTTP server write timeout (default: 30s)

##### Import OPML Command Options
- `--opml`, `-o`: Path to OPML file (required)
- `--db-path`: Path to the SQLite database file (default: ./riffle.db)

##### Run Command Options
- `--opml`, `-o`: Path to OPML file (required)
- `--interests`, `-i`: Path to file containing interests (one per line)
- `--articles`, `-n`: Number of articles to fetch from each feed (default: 3)
- `--top`, `-t`: Number of top articles to recommend (default: 1)
- `--model`, `-m`: Perplexity API model to use for article analysis (default: r1-1776)

## API Documentation

Riffle provides a comprehensive REST API for managing RSS sources, content, and recommendations. The API is documented in OpenAPI format.

For detailed API documentation, please refer to:
- [OpenAPI Specification](docs/api.yaml): Complete API specification in YAML format
- [API Documentation](docs/api_readme.md): Markdown version of the API documentation

The API includes endpoints for:
- RSS Source Management (CRUD operations)
- Content Management (fetching, updating, deleting)
- Content Search
- Recommendations
- System Information

## Frontend

The frontend provides a modern web interface for reading RSS feeds:

- Built with Vue.js and Material 3 design
- Displays all RSS sources on the left side of the page
- Shows the 10 most recent articles of each RSS source on the right side
- Automatically refreshes RSS content every 10 minutes

For more information about the frontend, see the [frontend README](frontend/README.md).

## License

This project is licensed under the Apache License, Version 2.0 - see the [LICENSE](LICENSE) file for details.