<template>
  <v-container fluid class="fill-height pa-0">
    <v-row no-gutters class="fill-height">
      <!-- Left panel: RSS Sources -->
      <v-col cols="3" class="border-right">
        <v-card flat class="fill-height rounded-0">
          <v-card-title class="d-flex align-center py-3">
            <span>RSS Sources</span>
            <v-spacer></v-spacer>
            <add-source-dialog @source-added="fetchSources" />
            <v-btn icon size="small" @click="refreshSources" class="ml-2">
              <v-icon>mdi-refresh</v-icon>
            </v-btn>
          </v-card-title>
          <v-divider></v-divider>
          <v-list v-if="sources && sources.length > 0">
            <v-list-item
              v-for="source in sources"
              :key="source.id"
              :title="source.name"
              :subtitle="source.description"
              :active="selectedSource && selectedSource.id === source.id"
              @click="selectSource(source)"
            >
              <template v-slot:prepend>
                <v-avatar color="primary" size="36">
                  <span class="text-h6 text-white">{{ source.name.charAt(0) }}</span>
                </v-avatar>
              </template>
            </v-list-item>
          </v-list>
          <div v-else-if="loading" class="d-flex justify-center align-center pa-4">
            <v-progress-circular indeterminate color="primary"></v-progress-circular>
          </div>
          <div v-else class="text-center pa-4">
            <v-icon size="64" color="grey-lighten-1">mdi-rss</v-icon>
            <div class="text-h6 mt-2">No RSS Sources</div>
            <div class="text-body-2 text-grey">Click the + button to add your first RSS source</div>
            <div class="mt-4">
              <sample-sources-dialog @sources-added="fetchSources" />
            </div>
          </div>
        </v-card>
      </v-col>

      <!-- Right panel: Articles -->
      <v-col cols="9">
        <v-card flat class="fill-height rounded-0">
          <v-card-title class="py-3">
            <span v-if="selectedSource">{{ selectedSource.name }}</span>
            <span v-else>Recommended Articles</span>
            <v-spacer></v-spacer>
            <v-btn icon size="small" @click="refreshArticles">
              <v-icon>mdi-refresh</v-icon>
            </v-btn>
          </v-card-title>
          <v-divider></v-divider>
          <v-card-text class="pa-0">
            <article-list 
              :articles="articles" 
              :loading="loading"
            ></article-list>
          </v-card-text>
        </v-card>
      </v-col>
    </v-row>
  </v-container>
</template>

<script>
import ApiService from '@/services/api'
import ArticleList from '@/components/ArticleList.vue'
import AddSourceDialog from '@/components/AddSourceDialog.vue'
import SampleSourcesDialog from '@/components/SampleSourcesDialog.vue'

export default {
  name: 'HomeView',
  components: {
    ArticleList,
    AddSourceDialog,
    SampleSourcesDialog
  },
  data() {
    return {
      sources: [],
      selectedSource: null,
      articles: [],
      loading: false,
      refreshInterval: null
    }
  },
  created() {
    this.fetchSources()
    this.fetchRecommendedArticles()
    
    // Set up auto-refresh every 10 minutes
    this.refreshInterval = setInterval(() => {
      this.refreshSources()
      this.refreshArticles()
    }, 10 * 60 * 1000)
    
    // Listen for refresh events
    this.emitter.on('refresh-feeds', this.handleRefresh)
  },
  beforeUnmount() {
    clearInterval(this.refreshInterval)
    this.emitter.off('refresh-feeds', this.handleRefresh)
  },
  methods: {
    async fetchSources() {
      try {
        this.loading = true
        const response = await ApiService.getSources()
        
        // Ensure we have a valid response with sources
        if (response && response.data && Array.isArray(response.data.sources)) {
          this.sources = response.data.sources
          console.log(`Loaded ${this.sources.length} sources`)
          
          // If we have a selected source, make sure it still exists
          if (this.selectedSource) {
            const sourceStillExists = this.sources.some(s => s.id === this.selectedSource.id)
            if (!sourceStillExists) {
              this.selectedSource = null
            }
          }
        } else {
          console.error('Invalid response format:', response)
          this.sources = []
        }
      } catch (error) {
        console.error('Error fetching sources:', error)
        this.sources = []
      } finally {
        this.loading = false
      }
    },
    async fetchRecommendedArticles() {
      try {
        this.loading = true
        const response = await ApiService.getRecommendations()
        this.articles = response.data.contents || []
        
        // If recommendations API fails or returns empty, get the latest articles
        if (this.articles.length === 0) {
          const latestResponse = await ApiService.getContents({ limit: 3 })
          this.articles = latestResponse.data.contents || []
        }
      } catch (error) {
        console.error('Error fetching recommended articles:', error)
        // Fallback to latest articles
        try {
          const latestResponse = await ApiService.getContents({ limit: 3 })
          this.articles = latestResponse.data.contents || []
        } catch (fallbackError) {
          console.error('Error fetching fallback articles:', fallbackError)
        }
      } finally {
        this.loading = false
      }
    },
    async fetchArticlesBySource(sourceId) {
      try {
        this.loading = true
        const response = await ApiService.getContentsBySource(sourceId, 10)
        this.articles = response.data.contents || []
      } catch (error) {
        console.error(`Error fetching articles for source ${sourceId}:`, error)
      } finally {
        this.loading = false
      }
    },
    selectSource(source) {
      this.selectedSource = source
      this.fetchArticlesBySource(source.id)
    },
    refreshSources() {
      this.fetchSources()
    },
    refreshArticles() {
      if (this.selectedSource) {
        this.fetchArticlesBySource(this.selectedSource.id)
      } else {
        this.fetchRecommendedArticles()
      }
    },
    handleRefresh() {
      this.refreshSources()
      this.refreshArticles()
    }
  }
}
</script>

<style scoped>
.border-right {
  border-right: 1px solid rgba(0, 0, 0, 0.12);
}
</style> 