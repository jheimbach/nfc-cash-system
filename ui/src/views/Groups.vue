<template>
  <div class="group-overview">
    <v-toolbar class="account-toolbar">
      <v-toolbar-title>Groups</v-toolbar-title>
      <v-spacer></v-spacer>
      <v-menu offset-y>
        <template v-slot:activator="{ on }">
          <v-btn
            color="success"
            dark
            v-on="on"
          >
            <v-icon left>group_add</v-icon>
            Add Group
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
    <group-list :groups="groups" :loading="loading" show-search/>
  </div>
</template>

<script lang="ts">

import { Component, Vue } from 'vue-property-decorator'
import Group from '../data/group'
import GroupList from '@/components/groups/List.vue'
import axios from 'axios'

@Component({
  components: { GroupList }
})
export default class Groups extends Vue {
  groups: Group[] = []
  loading: boolean = false
  addOptions: any = [
    {
      title: 'Add single Group',
      icon: 'add_box',
      to: { name: 'group_create' }
    },
    {
      title: 'Add multiple Groups',
      icon: 'playlist_add',
      to: { name: 'group_create_multiple' }
    },
    {
      title: 'Upload CSV File',
      icon: 'cloud_upload',
      to: { name: 'group_create_upload' }
    }
  ]

  getGroups() {
    axios.get('/groups', {
      baseURL: 'http://localhost:8088/v1',
      headers: {
        Authorization: 'Bearer ' + JSON.parse(localStorage.getItem('jwt')).access_token
      }
    }).then((response) => {
      this.groups = response.data.groups.map(el => {
        return {
          id: el.id,
          name: el.name,
          description: el.description,
          canOverdraw: el.can_overdraw
        }
      })
    }).catch((response) => {
      console.error(response)
    })
  }

  created() {
    this.loading = true
    setTimeout(() => {
      this.loading = false
      this.getGroups()
    }, 500)
  }
}
</script>

<style lang="scss" scoped>
  .group-overview {
    width: 100%;
  }
</style>
