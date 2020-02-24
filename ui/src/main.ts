import Vue from 'vue'
import App from './App.vue'
// @ts-ignore
import VueMaterial from 'vue-material'
import './registerServiceWorker'
import router from './router'

import formatDate from '@/filters/format-date'

import 'vue-material/dist/vue-material.min.css'
import './styles/main.scss'
import vuetify from './plugins/vuetify'

Vue.config.productionTip = false
Vue.use(VueMaterial)

Vue.filter('formatDate', formatDate)

new Vue({
  router: router,
  render: h => h(App),
  // @ts-ignore
  vuetify
}).$mount('#app')
