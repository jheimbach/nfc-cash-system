<template>
  <md-table v-model="searched" md-sort="name" md-sort-order="asc">
    <md-table-toolbar>
      <div class="md-toolbar-section-start">
        <h1 class="md-title">Accounts</h1>
      </div>

      <md-field md-clearable class="md-toolbar-section-end">
        <md-input placeholder="Search by name..." v-model="search" @input="searchOnTable"/>
      </md-field>
    </md-table-toolbar>

    <md-table-empty-state
      md-label="No accounts found"
      :md-description="`No accounts found for this '${search}' query. Try a different search term.`">
    </md-table-empty-state>

    <md-table-row slot="md-table-row" slot-scope="{ item }">
      <md-table-cell md-label="ID" md-sort-by="id" md-numeric>{{ item.id }}</md-table-cell>
      <md-table-cell md-label="Name" md-sort-by="name">{{ item.name }}</md-table-cell>
      <md-table-cell md-label="Description" md-sort-by="description">{{ item.description }}</md-table-cell>
      <md-table-cell md-label="Saldo" md-sort-by="saldo">{{ item.saldo }}</md-table-cell>
      <md-table-cell md-label="Group" md-sort-by="group.name">{{ item.group.name }}</md-table-cell>
      <md-table-cell md-label="NFC Chip Id" md-sort-by="nfcChipId">{{ item.nfcChipId }}</md-table-cell>
      <md-table-cell md-label="Actions">
        <md-button :to="{name: 'account', params:{id:item.id}}" class="md-icon-button md-primary">
          <md-icon>account_box</md-icon>
        </md-button>
        <md-button :to="{name: 'account_edit', params:{id:item.id}}" class="md-icon-button md-primary">
          <md-icon>edit</md-icon>
        </md-button>
      </md-table-cell>
    </md-table-row>
  </md-table>
</template>

<script lang="ts">

import { Component, Prop, Vue } from 'vue-property-decorator'

@Component
export default class AccountList extends Vue {
  @Prop({
    required: true,
    default: () => {
      return []
    }
  })
  accounts!: Account[]
  searched: Account[] = []
  private search: string = ''

  searchOnTable() {
    if (this.search) {
      let searchTerm = this.search.toLowerCase()
      this.searched = this.accounts.filter((item) => {
        if (item.name) {
          return item.name.toLowerCase().includes(searchTerm)
        }
      })
    }
  }

  created() {
    this.searched = this.accounts
  }
}
</script>

<style scoped>

</style>
