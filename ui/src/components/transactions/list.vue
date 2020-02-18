<template>
  <div class="transaction-list">
    <md-table v-model="searched" :md-sort="tableSort" :md-sort-order="tableSortOrder" md-fixed-header class="transaction-table" :md-sort-fn="sortByDate">
      <md-table-toolbar>
        <div class="md-toolbar-section-start">
          <h1 class="md-title">Transactions</h1>
        </div>

        <md-field md-clearable class="md-toolbar-section-end">
          <md-input placeholder="Search by name..." v-model="search" @input="searchOnTable"/>
        </md-field>
      </md-table-toolbar>

      <md-table-empty-state
        md-label="No transactions found"
        :md-description="`No transactions found for this '${search}' query. Try a different search term.`">
      </md-table-empty-state>

      <md-table-row slot="md-table-row" slot-scope="{ item }">
        <md-table-cell md-label="ID" md-sort-by="id" md-numeric>{{ item.id }}</md-table-cell>
        <md-table-cell md-label="Old Saldo">{{ item.oldSaldo }}</md-table-cell>
        <md-table-cell md-label="New Saldo">{{ item.newSaldo }}</md-table-cell>
        <md-table-cell md-label="Amount">{{ item.amount }}</md-table-cell>
        <md-table-cell md-label="Account" md-sort-by="account.name">{{ item.account.name }}</md-table-cell>
        <md-table-cell md-label="Created" md-sort-by="created">{{ item.created | formatDate }}</md-table-cell>
        <md-table-cell md-label="Actions" v-if="actions">
          <md-button :to="{ name: 'account', params: { id: item.account.id }}"
                     class="md-icon-button ncs-secondary">
            <md-icon>account_box</md-icon>
          </md-button>
        </md-table-cell>
      </md-table-row>
    </md-table>
  </div>
</template>

<script lang="ts">
import { Component, Prop, Vue } from 'vue-property-decorator'
import { Transaction } from '@/data/transaction'

@Component
export default class TransactionList extends Vue {
  search: string = ''
  tableSort: string = 'created'
  tableSortOrder: string = 'desc'
  searched: Transaction[] = []
  @Prop()
  transactions!: Transaction[]
  @Prop()
  showActions: boolean | undefined

  get actions() {
    return !this.showActions
  }

  sortByDate(value:Transaction[]) {
    console.log(this)
    console.log(this.tableSortOrder)
    return value.sort((a, b) => {
      return a.created > b.created ? 1 : a.created < b.created ? -1 : 0
    })
  }

  searchOnTable() {
    this.searched = searchByName(this.transactions, this.search)
  }

  created() {
    this.searched = this.transactions
  }
}

function toLower(text: string) {
  return text.toString().toLowerCase()
}

function searchByName(items: Transaction[], term: string) {
  if (term) {
    return items.filter((item: Transaction) => {
      return toLower(item.account.name).includes(toLower(term))
    })
  }

  return items
}
</script>

<style scoped>
</style>
