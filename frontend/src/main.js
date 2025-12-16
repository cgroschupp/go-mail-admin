import { createApp, markRaw } from 'vue'
import { createPinia } from 'pinia'
import App from './App.vue'
import router from './router'
import vuetify from './plugins/vuetify'
import { loadFonts } from './plugins/webfontloader'
import {Vuetify3Dialog} from 'vuetify3-dialog'
import piniaPluginPersistedstate from 'pinia-plugin-persistedstate'
import api, { setCsrfToken } from './service/Api'

const pinia = createPinia()
pinia.use(piniaPluginPersistedstate)
pinia.use(({ store }) => {
  store.router = markRaw(router)
})

loadFonts()

async function bootstrap() {
  const app = createApp(App)
  app.use(pinia)

  const { data } = await api().get('/api/v1/csrf')
  setCsrfToken(data.csrfToken)

  app.use(router)
    .use(vuetify)
    .use(Vuetify3Dialog, { vuetify })
    .mount('#app')
}

bootstrap()