# Riffle

Riffle is an RSS feed analyzer and content recommender that helps you find the most valuable content from your RSS subscriptions. It analyzes articles based on content quality and your personal interests to provide intelligent recommendations.

## Screenshot

![Riffle in action](docs/images/effect.png)

Example output showing article recommendations with interest matching and content quality scores.

## Features

- OPML feed list support
- Content quality analysis
- Personal interest matching
- Configurable article fetching
- Multiple recommendation levels
- Detailed scoring system

## Installation

To build the application:

```bash
make
```

This will create the binary in `bin/<os>_<arch>/riffle`.

## Usage

Basic usage:

```bash
riffle -o <opml-file> [-i <interests-file>] [-n <article-count>] [-t <top-count>]
```

### Command Line Options

| Flag | Short | Description | Default |
|------|-------|-------------|---------|
| `--opml` | `-o` | Path to OPML file (required) | - |
| `--interests` | `-i` | Path to file containing interests (one per line) | - |
| `--articles` | `-n` | Number of articles to fetch from each feed | 3 |
| `--top` | `-t` | Number of top articles to recommend | 1 |

### Examples

1. Basic usage with default settings (3 articles per feed, top 1 recommendation):
   ```bash
   riffle -o feeds.opml
   ```

2. Include personal interests for better recommendations:
   ```bash
   riffle -o feeds.opml -i interests.txt
   ```

3. Fetch more articles and get more recommendations:
   ```bash
   riffle -o feeds.opml -i interests.txt -n 5 -t 3
   ```

### Input File Formats

#### OPML File Format
The OPML file should contain your RSS feed subscriptions. Example:
```xml
<?xml version="1.0" encoding="UTF-8"?>
<opml version="1.0">
    <head>
        <title>My Feed Subscriptions</title>
    </head>
    <body>
        <outline text="Tech">
            <outline type="rss" 
                    text="Hacker News" 
                    title="Hacker News"
                    xmlUrl="https://news.ycombinator.com/rss" 
                    htmlUrl="https://news.ycombinator.com/"/>
        </outline>
    </body>
</opml>
```

#### Interests File Format
The interests file should contain one interest per line. Each interest can be multiple words. Example:
```
golang programming
web development
artificial intelligence
cloud computing
system architecture
performance optimization
```

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

## Development

### Prerequisites

- Go 1.24 or later
- Make

### Building for Different Platforms

To build for a specific platform:

```bash
make PLATFORMS=linux/amd64
```

Available platforms:
- linux/amd64
- linux/arm64
- darwin/amd64
- darwin/arm64

## License

This project is licensed under the Apache License 2.0 - see the [LICENSE](LICENSE) file for details. 