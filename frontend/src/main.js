import { createApp } from 'vue'
import mitt from 'mitt'
import App from './App.vue'
import router from './router'
import { createVuetify } from 'vuetify'
import * as components from 'vuetify/components'
import * as directives from 'vuetify/directives'
import 'vuetify/styles'
import { md3 } from 'vuetify/blueprints'

// Create a global event bus
const emitter = mitt()

const vuetify = createVuetify({
  components,
  directives,
  blueprint: md3,
  theme: {
    defaultTheme: 'light',
    themes: {
      light: {
        colors: {
          primary: '#006874',
          secondary: '#5c6bc0',
          accent: '#82B1FF',
          error: '#FF5252',
          info: '#2196F3',
          success: '#4CAF50',
          warning: '#FFC107'
        }
      }
    }
  }
})

const app = createApp(App)
// Add the event bus to the app's global properties
app.config.globalProperties.emitter = emitter
app.use(router)
app.use(vuetify)
app.mount('#app') 