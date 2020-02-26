import Vue from 'vue'
import App from './App.vue'
import './registerServiceWorker'
import router from './router'

import formatDate from '@/filters/format-date'
import formatMoney from '@/filters/format-money'

import vuetify from './plugins/vuetify'

Vue.config.productionTip = false

Vue.filter('formatDate', formatDate)
Vue.filter('formatMoney', formatMoney)

new Vue({
  router: router,
  render: h => h(App),
  // @ts-ignore
  vuetify
}).$mount('#app')
