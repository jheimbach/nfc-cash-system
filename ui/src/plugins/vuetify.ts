import Vue from 'vue'
// @ts-ignore
import Vuetify from 'vuetify/lib'
import de from 'vuetify/src/locale/de'
import en from 'vuetify/src/locale/en'

// @ts-ignore
import colors from 'vuetify/lib/util/colors'

Vue.use(Vuetify)

export default new Vuetify({
  lang: {
    locales: { de, en },
    current: 'en'
  },
  icons: {
    iconfont: 'md'
  },
  theme: {
    themes: {
      light: {
        primary: '#003b6b',
        secondary: '#f3712a'
      }
    }
  }
})
