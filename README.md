# Riffle

Riffle is an RSS feed analyzer and content recommender that helps you find the most valuable content from your RSS subscriptions. It analyzes articles based on content quality and your personal interests to provide intelligent recommendations.

## Features

- OPML feed list support
- Content quality analysis
- Personal interest matching
- Configurable article fetching
- Multiple recommendation levels
- Detailed scoring system
- Direct article URLs in output
- 2-day time window filtering

## Installation

To build the application:

```bash
make
```

This will create the binary in `bin/<os>_<arch>/riffle`.

### Docker

You can build and run Riffle using Docker in several ways:

#### 1. Build for Current Platform

```bash
# Build the image
make image

# Run with your configuration files
docker run --rm \
    -v "$(pwd)/conf/feeds.opml:/app/conf/feeds.opml:ro" \
    -v "$(pwd)/conf/interests.txt:/app/conf/interests.txt:ro" \
    -e OPENAI_API_KEY="your-api-key" \
    riffle:latest -o /app/conf/feeds.opml -i /app/conf/interests.txt
```

#### 2. Build for Specific Platform

```bash
# Build for Linux AMD64
make image-linux_amd64

# Build for Linux ARM64
make image-linux_arm64

# Build for macOS AMD64
make image-darwin_amd64

# Build for macOS ARM64
make image-darwin_arm64
```

#### Docker Run Options

You can customize the run configuration:

```bash
docker run --rm \
    -v "$(pwd)/conf/feeds.opml:/app/conf/feeds.opml:ro" \
    -v "$(pwd)/conf/interests.txt:/app/conf/interests.txt:ro" \
    -e OPENAI_API_KEY="your-api-key" \
    -e OPENAI_BASE_URL="your-api-base-url" \  # Optional
    riffle:latest \
    -o /app/conf/feeds.opml \
    -i /app/conf/interests.txt \
    -n 5 \                     # Number of articles to fetch
    -t 3 \                     # Number of recommendations
    -m "custom-model"          # Custom model name
```

## Environment Variables

The following environment variables are required for AI analysis functionality:

| Variable | Description | Required |
|----------|-------------|----------|
| `OPENAI_API_KEY` | API key for accessing the Perplexity AI service | Yes |
| `OPENAI_BASE_URL` | Base URL for the Perplexity API (defaults to https://api.perplexity.ai if not set) | No |

## Quick Start

1. Edit the configuration files to match your interests:
   - Add your RSS feeds to `feeds.opml` (see example in `conf/feeds.opml`)
   - Add your interests to `interests.txt` (see example in `conf/interests.txt`)

2. Run Riffle:
   ```bash
   export OPENAI_API_KEY="your-api-key"
   riffle -o feeds.opml -i interests.txt
   ```

## Usage

Basic usage:

```bash
riffle -o <opml-file> [-i <interests-file>] [-n <article-count>] [-t <top-count>] [-m <model-name>]
```

### Command Line Options

| Flag | Short | Description | Default |
|------|-------|-------------|---------|
| `--opml` | `-o` | Path to OPML file (required) | - |
| `--interests` | `-i` | Path to file containing interests (one per line) | - |
| `--articles` | `-n` | Number of articles to fetch from each feed (from last 2 days) | 3 |
| `--top` | `-t` | Number of most valuable articles to recommend | 1 |
| `--model` | `-m` | Perplexity API model to use for article analysis | r1-1776 |

### Output Format

The tool provides output in four main sections:

1. **Recent Articles** (📰)
   - Lists all articles from the last 2 days, grouped by RSS source
   - Shows title, URL, publication date, and summary (if available)
   - Only includes articles from your current time zone's last 48 hours

2. **Inactive Sources** (📅)
   - Lists RSS sources with no updates in the last 2 days
   - Helps track which feeds are currently inactive

3. **Most Valuable Articles** (🌟)
   - Unified recommendations across all sources
   - Shows detailed scoring and reasoning for each recommendation:
     * Interest Match Score (0.00-1.00)
     * Content Quality Score (0.00-1.00)
     * Overall Score (0.00-1.00)
     * Explanation of why the article was recommended

4. **AI Analysis** (📊)
   - Provides in-depth analysis of recommended articles using Perplexity AI

### Examples

1. Basic usage with default settings:
   ```bash
   riffle -o feeds.opml
   ```
   This will:
   - Process articles from the last 2 days
   - Fetch up to 3 articles per feed
   - Show 1 most valuable article recommendation
   - Use default AI model for analysis

2. Include personal interests for better recommendations:
   ```bash
   riffle -o feeds.opml -i interests.txt
   ```
   This will:
   - Match articles against your interests
   - Improve recommendation quality
   - Show interest match scores
   - Include AI analysis of top articles

3. Get more recommendations with custom model:
   ```bash
   riffle -o feeds.opml -i interests.txt -t 3 -m "custom-model"
   ```
   This will:
   - Show the top 3 most valuable articles
   - Rank them by overall score
   - Include detailed reasoning for each
   - Use specified AI model for analysis

4. Fetch more articles per feed:
   ```bash
   riffle -o feeds.opml -i interests.txt -n 5
   ```
   This will:
   - Fetch up to 5 articles per feed (from last 2 days)
   - Process more content for better selection
   - Still respect the 2-day time window
   - Include AI analysis of recommendations

Note: Articles are always filtered to the last 2 days in your current time zone, regardless of how many articles you request with `-n`. This ensures you only see recent, relevant content.

### Configuration Files

#### OPML File Format (feeds.opml)
The OPML file contains your RSS feed subscriptions, organized by category. Example structure:
```xml
<?xml version="1.0" encoding="UTF-8"?>
<opml version="1.0">
    <head>
        <title>Tech Feed Subscriptions</title>
    </head>
    <body>
        <outline text="Programming">
            <outline type="rss" 
                    text="The Go Blog" 
                    xmlUrl="https://go.dev/blog/feed.atom"/>
        </outline>
    </body>
</opml>
```

See `conf/feeds.opml` for a complete example with multiple feeds and categories.

#### Interests File Format (interests.txt)
The interests file contains your topics of interest, one per line. You can:
- Group interests with comments (lines starting with #)
- Use multiple words per interest
- Order by priority (all interests are weighted equally)

Example structure:
```
# Programming Languages
golang programming
rust development

# Technologies
cloud native
distributed systems
```

See `conf/interests.txt` for a complete example with various technology interests.

### Scoring System

Articles are scored based on multiple factors with equal weighting:

1. Content Quality (50%):
   - Text length (40% of quality score)
   - Keyword relevance (40% of quality score)
   - Link quality (20% of quality score)

2. Interest Match (50%):
   - Relevance to your specified interests
   - Matches are calculated using word-based analysis
   - Multiple word interests are handled intelligently

The final score is a combination of both factors, helping you find articles that are both high-quality and relevant to your interests.

## License

This project is licensed under the Apache License 2.0 - see the [LICENSE](LICENSE) file for details. 