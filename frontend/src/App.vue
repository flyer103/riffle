<template>
  <v-app>
    <v-app-bar color="primary" density="compact">
      <v-app-bar-title>Riffle RSS Reader</v-app-bar-title>
      <v-spacer></v-spacer>
      <v-btn icon @click="refreshFeeds">
        <v-icon>mdi-refresh</v-icon>
      </v-btn>
    </v-app-bar>

    <v-main>
      <router-view />
    </v-main>
  </v-app>
</template>

<script>
import ApiService from './services/api'

export default {
  name: 'App',
  methods: {
    async refreshFeeds() {
      try {
        await ApiService.fetchContents()
        this.emitter.emit('refresh-feeds')
      } catch (error) {
        console.error('Error refreshing feeds:', error)
      }
    }
  }
}
</script>

<style>
body {
  font-family: 'Roboto', sans-serif;
  margin: 0;
  padding: 0;
}

.v-application {
  background-color: #f5f5f5;
}

a {
  text-decoration: none;
}
</style> 