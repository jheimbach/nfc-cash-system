<template>
  <div class="account-overview">
    <v-toolbar class="account-toolbar">
      <v-toolbar-title>Accounts</v-toolbar-title>
      <v-spacer></v-spacer>
      <v-menu offset-y>
        <template v-slot:activator="{ on }">
          <v-btn
            color="success"
            dark
            v-on="on"
          >
            <v-icon left>person_add</v-icon>
            Add Account
          </v-btn>
        </template>
        <v-list>
          <v-list-item
            v-for="(item, index) in addOptions"
            :key="index"
            :to="item.to"
          >
            <v-list-item-title>
              <v-icon left>{{item.icon}}</v-icon>
              {{ item.title }}
            </v-list-item-title>
          </v-list-item>
        </v-list>
      </v-menu>
    </v-toolbar>
    <account-list :accounts="accounts" :loading="loading" show-actions show-search show-group/>
  </div>
</template>

<script lang="ts">

import { Component, Vue } from 'vue-property-decorator'
import Account from '@/data/account'
import AccountList from '@/components/account/List.vue'
import axios from 'axios'

@Component({
  components: { AccountList }
})

export default class Accounts extends Vue {
  accounts: Account[] = []
  loading: boolean = false
  addOptions: any = [
    {
      title: 'Add single Account',
      icon: 'add_box',
      to: { name: 'account_create' }
    },
    {
      title: 'Add multiple Accounts',
      icon: 'playlist_add',
      to: { name: 'account_create_multiple' }
    },
    {
      title: 'Upload CSV File',
      icon: 'cloud_upload',
      to: { name: 'account_create_upload' }
    }
  ]

  getAccounts() {
    axios.get('/accounts', {
      baseURL: 'http://localhost:8088/v1',
      headers: {
        Authorization: 'Bearer ' + JSON.parse(localStorage.getItem('jwt')).access_token
      }
    }).then((response) => {
      this.accounts = response.data.accounts
    }).catch((response) => {
      console.error(response)
    })
  }

  created() {
    this.loading = true
    setTimeout(() => {
      this.loading = false
      this.getAccounts()
    }, 500)
  }
}
</script>

<style lang="scss" scoped>
  .md-field {
    max-width: 300px;
  }
</style>

<style lang="scss">
  .account-overview {
    width: 100%;
  }

  .account-toolbar {
    margin-bottom: 20px;
  }
</style>
