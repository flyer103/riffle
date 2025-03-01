<template>
  <div>
    <v-dialog v-model="dialog" max-width="500px">
      <template v-slot:activator="{ props }">
        <v-btn
          color="secondary"
          v-bind="props"
          variant="text"
          class="ml-2"
        >
          Add Sample Sources
        </v-btn>
      </template>
      <v-card>
        <v-card-title>
          <span class="text-h6">Add Sample RSS Sources</span>
        </v-card-title>
        <v-card-text>
          <p>Add these sample RSS sources to get started:</p>
          <v-list>
            <v-list-item v-for="(source, index) in sampleSources" :key="index">
              <v-list-item-title>{{ source.name }}</v-list-item-title>
              <v-list-item-subtitle>{{ source.url }}</v-list-item-subtitle>
            </v-list-item>
          </v-list>
        </v-card-text>
        <v-card-actions>
          <v-spacer></v-spacer>
          <v-btn color="blue-darken-1" variant="text" @click="close">
            Cancel
          </v-btn>
          <v-btn color="blue-darken-1" variant="text" @click="addSamples" :loading="loading">
            Add All
          </v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>
  </div>
</template>

<script>
import ApiService from '@/services/api'

export default {
  name: 'SampleSourcesDialog',
  data() {
    return {
      dialog: false,
      loading: false,
      sampleSources: [
        {
          name: 'BBC News',
          url: 'http://feeds.bbci.co.uk/news/rss.xml',
          description: 'The latest stories from the BBC'
        },
        {
          name: 'CNN',
          url: 'http://rss.cnn.com/rss/edition.rss',
          description: 'CNN news feed'
        },
        {
          name: 'Hacker News',
          url: 'https://news.ycombinator.com/rss',
          description: 'Hacker News RSS feed'
        },
        {
          name: 'TechCrunch',
          url: 'https://techcrunch.com/feed/',
          description: 'TechCrunch RSS feed'
        }
      ]
    }
  },
  methods: {
    async addSamples() {
      try {
        this.loading = true
        
        // Use batch create to add all sources at once
        await ApiService.batchCreateSources(this.sampleSources)
        
        // Fetch content for all sources
        try {
          await ApiService.fetchContents()
        } catch (fetchError) {
          console.error('Error fetching content for sources:', fetchError)
          // Continue even if fetch fails
        }
        
        this.close()
        this.$emit('sources-added')
      } catch (error) {
        console.error('Error adding sample sources:', error)
        
        // Fallback to individual creation if batch fails
        try {
          for (const source of this.sampleSources) {
            try {
              const sourceResponse = await ApiService.createSource(source)
              // Try to fetch content for this source
              try {
                await ApiService.fetchContentForSource(sourceResponse.data.id)
              } catch (fetchError) {
                console.error(`Error fetching content for source ${source.name}:`, fetchError)
              }
            } catch (sourceError) {
              console.error(`Error adding sample source ${source.name}:`, sourceError)
              // Continue with other sources even if one fails
            }
          }
          this.close()
          this.$emit('sources-added')
        } catch (fallbackError) {
          console.error('Error in fallback source creation:', fallbackError)
        }
      } finally {
        this.loading = false
      }
    },
    close() {
      this.dialog = false
    }
  }
}
</script> 