# Riffle Frontend

This is the frontend for the Riffle RSS reader application. It provides a modern web interface for reading RSS feeds.

## Features

- Modern web interface using Vue.js and Material 3 design
- Display of all RSS sources on the left side of the page
- Display of the 10 most recent articles of each RSS source on the right side
- Automatic refresh of RSS content every 10 minutes

## Project Setup

```bash
# Install dependencies
npm install

# Serve with hot-reload for development
npm run serve

# Build for production
npm run build

# Lint and fix files
npm run lint
```

## Configuration

The frontend is configured to connect to the Riffle backend API running on `http://localhost:8080`. If your backend is running on a different URL, you can modify the `baseURL` in `src/services/api.js` or update the proxy settings in `vue.config.js`.

## Architecture

- **Vue.js**: Frontend framework
- **Vuetify**: Material 3 design components
- **Axios**: HTTP client for API requests
- **Vue Router**: Client-side routing

## Directory Structure

- `src/assets`: Static assets
- `src/components`: Vue components
- `src/views`: Vue views (pages)
- `src/services`: API services
- `src/router`: Vue Router configuration 