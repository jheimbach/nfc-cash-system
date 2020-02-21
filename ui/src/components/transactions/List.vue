<template>
  <div class="transaction-list">
    <md-table v-model="searched" :md-sort="tableSort" :md-sort-order="tableSortOrder" md-fixed-header
              class="transaction-table" :md-sort-fn="sortByDate" :md-height="665">
      <md-table-toolbar>
        <div class="md-toolbar-section-start">
          <h1 class="md-title">Transactions</h1>
        </div>

        <md-field md-clearable class="md-toolbar-section-end" v-if="showSearch">
          <md-input placeholder="Search..." v-model="search" @input="searchOnTable"/>
        </md-field>
      </md-table-toolbar>

      <md-table-empty-state
        md-icon="account_balance_wallet"
        md-label="No transactions found"
        :md-description="emptyMsg">
      </md-table-empty-state>

      <md-table-row slot="md-table-row" slot-scope="{ item }">
        <md-table-cell md-label="ID" md-sort-by="id" md-numeric>{{ item.id }}</md-table-cell>
        <md-table-cell md-label="Old Saldo" md-numeric>{{ item.oldSaldo }}</md-table-cell>
        <md-table-cell md-label="New Saldo" md-numeric>{{ item.newSaldo }}</md-table-cell>
        <md-table-cell md-label="Amount" md-numeric>{{ item.amount }}</md-table-cell>
        <md-table-cell md-label="Account" md-sort-by="account.name">{{ item.account.name }}</md-table-cell>
        <md-table-cell md-label="Created" md-sort-by="created">{{ item.created | formatDate }}</md-table-cell>
        <md-table-cell md-label="Actions" v-if="showActions">
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

  @Prop({
    default: () => {
      return []
    },
    type: Array,
    required: true
  })
  transactions!: Transaction[]

  @Prop({
    default: false,
    type: Boolean
  })
  showActions!: boolean

  @Prop({
    default: false,
    type: Boolean
  })
  showSearch!: boolean

  @Prop({
    default: `No transactions found`,
    type: String
  })
  private emptyMessage!: string

  get emptyMsg(): string {
    if (this.search !== '') {
      return `No transactions found for this '${this.search}' query. Try a different search term.`
    }

    return this.emptyMessage
  }

  sortByDate(value: Transaction[]) {
    return value.sort((a, b) => {
      const { MdTable } = this.$children[0] as any
      const sortBy = MdTable.sort
      let aAttr = getObjectAttribute(a, sortBy)
      let bAttr = getObjectAttribute(b, sortBy)
      const isAsc = MdTable.sortOrder === 'asc'
      let isNumber = typeof aAttr === 'number'
      let isDate = sortBy === 'created'

      if (!aAttr) {
        return 1
      }

      if (!bAttr) {
        return -1
      }

      if (isDate) {
        if (isAsc) {
          return (aAttr < bAttr) ? -1 : bAttr < aAttr ? 1 : 0
        } else {
          return (aAttr < bAttr) ? 1 : bAttr < aAttr ? -1 : 0
        }
      }

      if (isNumber) {
        return isAsc ? (aAttr - bAttr) : (bAttr - aAttr)
      }

      return isAsc
        ? aAttr.localeCompare(bAttr)
        : bAttr.localeCompare(aAttr)
    })
  }

  searchOnTable() {
    this.searched = searchByName(this.transactions, this.search)
  }

  created() {
    this.searched = this.transactions
  }
}

function getObjectAttribute(object: Transaction, key: string) {
  let value: any = object

  for (const attribute of key.split('.')) {
    value = value[attribute]
  }

  return value
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

<style scoped lang="scss">
  .md-toolbar {
    padding-left: 0;

    .md-title {
      margin-left: 0;
    }
  }

  .md-table-cell {
    padding-left: 0;
  }
</style>
