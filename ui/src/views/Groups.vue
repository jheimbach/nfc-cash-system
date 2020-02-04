<template>
  <div>
    <md-table v-model="searched" md-sort="name" md-sort-order="asc">
      <md-table-toolbar>
        <div class="md-toolbar-section-start">
          <h1 class="md-title">Groups</h1>
        </div>

        <md-field md-clearable class="md-toolbar-section-end">
          <md-input placeholder="Search by name..." v-model="search" @input="searchOnTable"/>
        </md-field>
      </md-table-toolbar>

      <md-table-empty-state
        md-label="No groups found"
        :md-description="`No groups found for this '${search}' query. Try a different search term.`">
      </md-table-empty-state>

      <md-table-row slot="md-table-row" slot-scope="{ item }">
        <md-table-cell md-label="ID" md-sort-by="id" md-numeric>{{ item.id }}</md-table-cell>
        <md-table-cell md-label="Name" md-sort-by="name">{{ item.name }}</md-table-cell>
        <md-table-cell md-label="Descripion" md-sort-by="description">{{ item.description }}</md-table-cell>
        <md-table-cell md-label="Can overdraw" md-sort-by="can_overdraw">{{ item.can_overdraw ? 'Yes' :'No' }}
        </md-table-cell>
        <md-table-cell md-label="Actions">
          <md-button :to="accountsRoute(item.id)" class="md-icon-button ncs-secondary">
            <md-icon>account_box</md-icon>
          </md-button>
          <md-button :to="editRoute(item.id)" class="md-icon-button md-primary">
            <md-icon>edit</md-icon>
          </md-button>
          <md-button :to="deleteRoute(item.id)" class="md-icon-button md-accent">
            <md-icon>delete</md-icon>
          </md-button>
        </md-table-cell>
      </md-table-row>
    </md-table>
  </div>
</template>

<script lang="ts">

import { Component, Vue } from 'vue-property-decorator'
import Group from '../data/group'

@Component
export default class Groups extends Vue {
  search?: string
  searched: Group[] = []
  groups: Group[] = []

  searchOnTable() {
    if (this.search) {
      let searchTerm = this.search.toLowerCase()
      this.searched = this.groups.filter((item) => item.name.toLowerCase().includes(searchTerm))
    }
  }

  editRoute(groupId: number) {
    return `/groups/${groupId}/edit`
  }

  accountsRoute(groupId: number) {
    return `/groups/${groupId}/accounts`
  }

  deleteRoute(groupId: number) {
    return `/groups/${groupId}/delete`
  }

  created() {
    this.searched = this.groups
  }
}
</script>

<style lang="scss" scoped>
  .md-field {
    max-width: 300px;
  }
</style>
