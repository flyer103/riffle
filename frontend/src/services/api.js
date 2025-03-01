import axios from 'axios'

const apiClient = axios.create({
  baseURL: 'http://localhost:8080',
  withCredentials: false,
  headers: {
    Accept: 'application/json',
    'Content-Type': 'application/json'
  },
  timeout: 10000
})

export default {
  // RSS Sources
  getSources() {
    return apiClient.get('/sources')
  },
  getSource(id) {
    return apiClient.get(`/sources/${id}`)
  },
  createSource(source) {
    return apiClient.post('/sources', source)
  },
  batchCreateSources(sources) {
    return apiClient.post('/sources/batch', { sources })
  },
  updateSource(id, source) {
    return apiClient.put(`/sources/${id}`, source)
  },
  deleteSource(id) {
    return apiClient.delete(`/sources/${id}`)
  },

  // RSS Contents
  getContents(params = {}) {
    return apiClient.get('/contents', { params })
  },
  getContent(id) {
    return apiClient.get(`/contents/${id}`)
  },
  getContentsBySource(sourceId, limit = 10) {
    return apiClient.get('/contents', { 
      params: { 
        sourceId, 
        limit 
      } 
    })
  },
  
  // Recommendations
  getRecommendations() {
    return apiClient.get('/recommendations')
  },

  // Fetch new content
  fetchContents() {
    return apiClient.post('/contents/fetch', { days: 30 })
  },
  fetchContentForSource(sourceId) {
    return apiClient.post('/contents/fetch', { 
      sourceId, 
      days: 30 
    })
  }
} 