# Riffle API Documentation

This directory contains the OpenAPI specification for the Riffle API.

## OpenAPI Specification

The API is documented using the [OpenAPI 3.0](https://swagger.io/specification/) specification format in the `api.yaml` file. This is a standardized, language-agnostic interface description for HTTP APIs that allows both humans and computers to understand the capabilities of a service without direct access to source code or documentation.

## Using the API Documentation

### Viewing the Documentation

You can view the API documentation in several ways:

1. **Swagger UI**: Import the `api.yaml` file into [Swagger Editor](https://editor.swagger.io/) to view and interact with the API documentation.

2. **Redoc**: Use [Redoc](https://redocly.github.io/redoc/) to view a more user-friendly version of the documentation.

3. **Postman**: Import the `api.yaml` file into [Postman](https://www.postman.com/) to create a collection of API requests that you can use to test the API.

### Generating Client Libraries

You can use the OpenAPI specification to generate client libraries for various programming languages:

1. **OpenAPI Generator**: Use [OpenAPI Generator](https://openapi-generator.tech/) to generate client libraries for over 40 different languages.

```bash
# Install OpenAPI Generator
npm install @openapitools/openapi-generator-cli -g

# Generate a client library (example for Python)
openapi-generator-cli generate -i api.yaml -g python -o ./python-client
```

2. **Swagger Codegen**: Use [Swagger Codegen](https://github.com/swagger-api/swagger-codegen) as an alternative.

## API Overview

The Riffle API provides endpoints for:

1. **RSS Sources Management**: Create, read, update, and delete RSS sources
2. **RSS Contents Management**: Retrieve, update, and delete content items
3. **Content Fetching**: Fetch new content from RSS sources
4. **Recommendations**: Get content recommendations and submit user feedback
5. **System Information**: Check system health and get system information

## Using with the import-opml Command

The `import-opml` command can be used to import RSS sources from an OPML file into the database:

```bash
riffle import-opml -o path/to/opml/file.opml --db-path ./riffle.db
```

After importing sources, you can use the API to:

1. List the imported sources: `GET /sources`
2. Check the status of the fetch job: `GET /contents/fetch/{jobId}`
3. Get content from the imported sources: `GET /contents?sourceId={sourceId}`
4. Get recommendations based on the imported content: `GET /recommendations`

## Example Workflow

1. Import sources from an OPML file using the `import-opml` command
2. Start the server using the `serve` command
3. Use the API to manage sources, fetch content, and get recommendations 