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
            <v-icon left>add_box</v-icon>
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
    this.groups = [
      {
        id: 1,
        name: 'H2O Plus'
      },
      {
        id: 2,
        name: 'A-S Medication Solutions LLC',
        description: 'E.E.S',
        canOverdraw: true
      },
      {
        id: 3,
        name: 'Mylan Pharmaceuticals Inc.'
      },
      {
        id: 4,
        name: 'Mylan Pharmaceuticals Inc.',
        description: 'Enalapril Maleate and Hydrochlorothiazide'
      },
      {
        id: 5,
        name: 'REMEDYREPACK INC.',
        description: 'CELEBREX'
      },
      {
        id: 6,
        name: 'H E B',
        description: 'night time'
      },
      {
        id: 7,
        name: 'PSS World Medical, Inc.',
        canOverdraw: true
      },
      {
        id: 8,
        name: 'Kareway Product, Inc.',
        description: 'Acetaminophen',
        canOverdraw: true
      },
      {
        id: 9,
        name: 'Pharmacia and Upjohn Company'
      },
      {
        id: 10,
        name: 'Dolgencorp, Inc. (DOLLAR GENERAL & REXALL)',
        description: 'Allergy Relief',
        canOverdraw: true
      }
    ]
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
