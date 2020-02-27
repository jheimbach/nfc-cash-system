<template>
  <v-app>
    <v-content>
      <v-container
        class="fill-height"
        fluid
      >
        <v-row align="center" justify="center">
          <v-col
            cols="12"
            sm="8"
            md="4"
            class="login-box"
          >
            <v-card>
              <v-toolbar
                color="primary"
                dark
                flat
              >
                <v-toolbar-title>Login form</v-toolbar-title>
              </v-toolbar>
              <v-card-text>
                <v-form>
                  <v-text-field
                    v-model="email"
                    label="Login"
                    prepend-icon="person"
                    type="text"
                  />

                  <v-text-field
                    v-model="password"
                    label="Password"
                    prepend-icon="lock"
                    type="password"
                  />
                </v-form>
              </v-card-text>
              <v-card-actions>
                <v-spacer/>
                <v-btn color="secondary" @click="handleSubmit">Login</v-btn>
              </v-card-actions>
            </v-card>
          </v-col>
        </v-row>
      </v-container>
    </v-content>
  </v-app>
</template>

<script lang="ts">
import { Component, Vue } from 'vue-property-decorator'
import axios from 'axios'

@Component
export default class Login extends Vue {
  email: string = 'briste0@angelfire.com'
  password: string = 'lMvZARjM3pwe'

  handleSubmit(e: MouseEvent) {
    e.preventDefault()
    if (this.password.length > 0) {
      axios.get('/user/login', {
        baseURL: 'http://localhost:8088/v1',
        auth: {
          username: this.email,
          password: this.password
        },
        headers: {
          'Content-type': 'application/json',
          'Accept': 'application/json'
        }
      })
        .then(response => {
          localStorage.setItem('jwt', JSON.stringify(response.data))
          const route = this.$route.query.to ? JSON.parse(this.$route.query.to.toString()) : { name: 'home' }
          this.$router.push(route)
        })
        .catch(function (error) {
          console.error(error)
        })
    }
  }
}
</script>
<style scoped lang="scss">
  .container.fill-height {
    > .row {
      max-width: calc(100% + 24px);
    }
  }
</style>
