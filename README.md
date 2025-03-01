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

## Getting Started

### Prerequisites

- Go 1.24 or higher
- SQLite3

### Installation

1. Clone the repository:

```bash
git clone https://github.com/flyer103/riffle.git
cd riffle
```

2. Build the application:

```bash
make build
```

### Usage

Riffle provides several commands:

#### Running the Server

```bash
./riffle serve --port 8080 --db-path ./riffle.db
```

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

### RSS Sources

#### List Sources

```
GET /sources
```

Query Parameters:
- `limit`: Maximum number of sources to return (default: 50)
- `nextToken`: Token for pagination

#### Get Source

```
GET /sources/:id
```

#### Create Source

```
POST /sources
```

Request Body:
```json
{
  "name": "Example Feed",
  "url": "https://example.com/feed.xml",
  "description": "An example RSS feed",
  "category": "Technology"
}
```

#### Update Source

```
PUT /sources/:id
```

Request Body:
```json
{
  "name": "Updated Feed Name",
  "description": "Updated description",
  "category": "News"
}
```

#### Delete Source

```
DELETE /sources/:id
```

#### Batch Create Sources

```
POST /sources/batch
```

Request Body:
```json
{
  "sources": [
    {
      "name": "Feed 1",
      "url": "https://example1.com/feed.xml",
      "category": "Technology"
    },
    {
      "name": "Feed 2",
      "url": "https://example2.com/feed.xml",
      "category": "News"
    }
  ]
}
```

#### Batch Delete Sources

```
DELETE /sources/batch
```

Request Body:
```json
{
  "sourceIds": ["source-id-1", "source-id-2"]
}
```

### RSS Contents

#### List Contents

```
GET /contents
```

Query Parameters:
- `sourceId`: Filter by source ID
- `startDate`: Filter by start date (RFC3339 format)
- `endDate`: Filter by end date (RFC3339 format)
- `limit`: Maximum number of contents to return (default: 50)
- `nextToken`: Token for pagination

#### Get Content

```
GET /contents/:id
```

#### Update Content

```
PUT /contents/:id
```

Request Body:
```json
{
  "title": "Updated Title",
  "description": "Updated description",
  "content": "Updated content",
  "categories": ["Technology", "Programming"]
}
```

#### Delete Content

```
DELETE /contents/:id
```

#### Batch Delete Contents

```
DELETE /contents/batch
```

Request Body:
```json
{
  "contentIds": ["content-id-1", "content-id-2"]
}
```

#### Fetch Contents

```
POST /contents/fetch
```

Request Body:
```json
{
  "sourceId": "source-id",  // Optional, if not provided, fetch all sources
  "days": 7                 // Number of days to fetch
}
```

#### Get Fetch Status

```
GET /contents/fetch/:jobId
```

#### Search Contents

```
GET /contents/search?keywords=golang,programming&sourceId=source-id&limit=10
```

Query Parameters:
- `keywords`: Comma-separated list of keywords to search for
- `sourceId`: Filter by source ID (optional)
- `limit`: Maximum number of results to return (default: 50)

### Recommendations

#### Get Recommendations

```
GET /recommendations
```

Query Parameters:
- `userId`: User ID to get recommendations for
- `sourceIds`: Comma-separated list of source IDs to filter by
- `limit`: Maximum number of recommendations to return (default: 10)

#### Submit Feedback

```
POST /recommendations/feedback
```

Request Body:
```json
{
  "contentId": "content-id",
  "userId": "user-id",
  "rating": 5,
  "comment": "Great article!"
}
```

#### Get User Feedback

```
GET /recommendations/feedback/:userId
```

### System

#### Health Check

```
GET /health
```

#### System Info

```
GET /system/info
```

## License

This project is licensed under the MIT License - see the LICENSE file for details.