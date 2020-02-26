import Vue from 'vue'
// @ts-ignore
import Vuetify from 'vuetify/lib/framework'
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
    options: {
      customProperties: true
    },
    themes: {
      light: {
        primary: '#003b6b',
        secondary: '#f3712a',
        accent: colors.teal.base,
        error: colors.red.base,
        warning: colors.amber.base,
        info: colors.blueGrey.base,
        success: colors.lightGreen.base
      }
    }
  }
})
