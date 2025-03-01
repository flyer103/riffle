<template>
  <div>
    <v-progress-linear
      v-if="loading"
      indeterminate
      color="primary"
    ></v-progress-linear>
    
    <div v-else-if="articles.length === 0" class="text-center pa-8">
      <v-icon size="64" color="grey lighten-1">mdi-rss</v-icon>
      <div class="text-h6 mt-4 text-grey">No articles found</div>
    </div>
    
    <v-list v-else lines="three">
      <v-list-item
        v-for="article in articles"
        :key="article.id"
        :title="article.title"
        :subtitle="formatDate(article.publishedAt)"
      >
        <template v-slot:prepend>
          <v-avatar color="grey-lighten-3" size="48" class="mr-3">
            <v-icon color="primary">mdi-newspaper</v-icon>
          </v-avatar>
        </template>
        
        <template v-slot:append>
          <v-btn
            variant="text"
            color="primary"
            :href="article.link"
            target="_blank"
            rel="noopener noreferrer"
          >
            Read
          </v-btn>
        </template>
        
        <v-list-item-subtitle class="mt-2 text-truncate-3">
          {{ stripHtml(article.description) }}
        </v-list-item-subtitle>
      </v-list-item>
    </v-list>
  </div>
</template>

<script>
export default {
  name: 'ArticleList',
  props: {
    articles: {
      type: Array,
      default: () => []
    },
    loading: {
      type: Boolean,
      default: false
    }
  },
  methods: {
    formatDate(dateString) {
      if (!dateString) return ''
      
      const date = new Date(dateString)
      return new Intl.DateTimeFormat('en-US', {
        year: 'numeric',
        month: 'short',
        day: 'numeric',
        hour: '2-digit',
        minute: '2-digit'
      }).format(date)
    },
    stripHtml(html) {
      if (!html) return ''
      
      // Create a temporary element to strip HTML tags
      const temp = document.createElement('div')
      temp.innerHTML = html
      return temp.textContent || temp.innerText || ''
    }
  }
}
</script>

<style scoped>
.text-truncate-3 {
  display: -webkit-box;
  -webkit-line-clamp: 3;
  -webkit-box-orient: vertical;
  overflow: hidden;
  text-overflow: ellipsis;
}
</style> 