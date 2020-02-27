<template>
  <v-card width="100%">
    <v-card-title>
      Transactions
      <v-spacer/>
      <v-text-field v-model="search" append-icon="search" label="Search" single-line hide-details v-if="showSearch"/>
    </v-card-title>
    <v-data-table :headers="headers" :items="searched" :items-per-page="itemsPerPage" fixed-header item-key="id"
                  :search="search" sort-by="created" sort-desc
                  :loading="loading" :hide-default-footer="searched.length < itemsPerPage">
      <template v-slot:item.action="{ item }" v-if="showActions">
        <v-icon small class="mr-2" @click="editItem(item)">account_box</v-icon>
      </template>
      <template v-slot:item.created="{ item }">
        {{item.created|formatDate}}
      </template>
      <template v-slot:item.amount="{ item }">
        {{item.amount|formatMoney}}
      </template>
      <template v-slot:item.oldSaldo="{ item }">
        {{item.oldSaldo|formatMoney}}
      </template>
      <template v-slot:item.newSaldo="{ item }">
        {{item.newSaldo|formatMoney}}
      </template>
    </v-data-table>
  </v-card>
</template>

<script lang="ts">
import { Component, Prop, Vue, Watch } from 'vue-property-decorator'
import Transaction from '@/data/transaction'
import { DataTableHeader } from 'vuetify'

@Component
export default class TransactionList extends Vue {
  search: string = ''
  searched: Transaction[] = []
  itemsPerPage: number = 25
  defaultHeaders: DataTableHeader[] = [
    {
      text: 'ID',
      value: 'id',
      align: 'end'
    },
    {
      text: 'Amount',
      value: 'amount',
      align: 'end'
    },
    {
      text: 'Old Saldo',
      value: 'oldSaldo',
      align: 'end'
    },
    {
      text: 'New Saldo',
      value: 'newSaldo',
      align: 'end'
    },
    {
      text: 'Created',
      value: 'created'
    }
  ]

  get headers(): DataTableHeader[] {
    if (this.showAccountName) {
      this.defaultHeaders.splice(2, 0, {
        text: 'Account',
        value: 'account.name'
      })
    }
    if (this.showActions) {
      // push as last element the actions column
      this.defaultHeaders.push(
        {
          text: '',
          value: 'action',
          sortable: false
        })
    }
    return this.defaultHeaders
  }

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
    type: Boolean,
    required: false
  })
  loading!: boolean

  @Prop({
    default: false,
    type: Boolean
  })
  showActions!: boolean

  @Prop({
    default: false,
    type: Boolean
  })
  showAccountName!: boolean

  @Prop({
    default: false,
    type: Boolean
  })
  showSearch!: boolean

  @Watch('transactions')
  updateSearched() {
    this.searched = this.transactions
  }

  editItem(item: Transaction) {
    this.$router.push({ name: 'account', params: { id: item.account.id.toString() } })
  }

  created() {
    this.searched = this.transactions
  }
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
