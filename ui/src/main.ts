import Vue from 'vue'
import App from './App.vue'
// @ts-ignore
import VueMaterial from 'vue-material'
import './registerServiceWorker'
import router from './router'

import formatDate from '@/filters/format-date'

import 'vue-material/dist/vue-material.min.css'
import 'vue-material/dist/theme/default.css'

Vue.config.productionTip = false
Vue.use(VueMaterial)

Vue.filter('formatDate', formatDate)

const linkActiveClass: string = 'ncs-active-class'
// @ts-ignore
Vue.material.router.linkActiveClass = linkActiveClass

new Vue({
  router: router(linkActiveClass),
  render: h => h(App)
}).$mount('#app')
