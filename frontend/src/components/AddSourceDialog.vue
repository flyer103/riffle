<template>
  <div>
    <v-dialog v-model="dialog" max-width="500px">
      <template v-slot:activator="{ props }">
        <v-btn
          color="primary"
          v-bind="props"
          size="small"
          icon
        >
          <v-icon>mdi-plus</v-icon>
        </v-btn>
      </template>
      <v-card>
        <v-card-title>
          <span class="text-h6">Add RSS Source</span>
        </v-card-title>
        <v-card-text>
          <v-container>
            <v-row>
              <v-col cols="12">
                <v-text-field
                  v-model="source.name"
                  label="Name"
                  required
                ></v-text-field>
              </v-col>
              <v-col cols="12">
                <v-text-field
                  v-model="source.url"
                  label="URL"
                  required
                ></v-text-field>
              </v-col>
              <v-col cols="12">
                <v-textarea
                  v-model="source.description"
                  label="Description"
                  rows="3"
                ></v-textarea>
              </v-col>
            </v-row>
          </v-container>
        </v-card-text>
        <v-card-actions>
          <v-spacer></v-spacer>
          <v-btn color="blue-darken-1" variant="text" @click="close">
            Cancel
          </v-btn>
          <v-btn color="blue-darken-1" variant="text" @click="save" :loading="loading">
            Save
          </v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>
  </div>
</template>

<script>
import ApiService from '@/services/api'

export default {
  name: 'AddSourceDialog',
  data() {
    return {
      dialog: false,
      loading: false,
      source: {
        name: '',
        url: '',
        description: ''
      }
    }
  },
  methods: {
    async save() {
      if (!this.source.name || !this.source.url) {
        alert('Name and URL are required')
        return
      }
      
      try {
        this.loading = true
        const response = await ApiService.createSource(this.source)
        
        // Fetch content for the new source
        try {
          await ApiService.fetchContentForSource(response.data.id)
        } catch (fetchError) {
          console.error('Error fetching content for new source:', fetchError)
          // Continue even if fetch fails
        }
        
        this.close()
        this.$emit('source-added')
      } catch (error) {
        console.error('Error creating source:', error)
        alert('Failed to create source: ' + (error.response?.data?.error || error.message))
      } finally {
        this.loading = false
      }
    },
    close() {
      this.dialog = false
      this.source = {
        name: '',
        url: '',
        description: ''
      }
    }
  }
}
</script> 