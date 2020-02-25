<template>
  <v-card width="100%">
    <v-card-title>
      <template v-if="showTitle">Accounts</template>
      <v-spacer v-if="showTitle && showSearch"/>
      <v-text-field v-model="search" append-icon="search" label="Search" single-line hide-details v-if="showSearch"/>
    </v-card-title>
    <v-data-table :headers="headers" :items="searched" :items-per-page="itemsPerPage" fixed-header item-key="id"
                  :search="search" sort-by="id"
                  :loading="loading" :hide-default-footer="searched.length < itemsPerPage">
      <template v-slot:item.action="{ item }">
        <v-icon small class="mr-2" @click="editItem(item)">edit</v-icon>
        <v-icon small @click="deleteItem(item)">delete</v-icon>
      </template>
    </v-data-table>
  </v-card>
</template>

<script lang="ts">

import { Component, Emit, Prop, Vue, Watch } from 'vue-property-decorator'
import Account from '@/data/account'
import { DataTableHeader } from 'vuetify'

@Component
export default class AccountList extends Vue {
  defaultHeaders: DataTableHeader[] = [
    {
      text: 'ID',
      align: 'end',
      value: 'id',
      filterable: false
    },
    {
      text: 'Name',
      value: 'name'
    },
    {
      text: 'Description',
      value: 'description',
      sortable: false
    },
    {
      text: 'Saldo',
      align: 'end',
      value: 'saldo',
      filterable: false
    },
    {
      text: 'NFC Chip ID',
      align: 'end',
      value: 'nfcChipId',
      sortable: false
    }
  ]

  get headers() {
    if (this.showGroup) {
      this.defaultHeaders.splice(3, 0,
        {
          text: 'Group',
          value: 'group.name',
          sortable: false
        })
    }
    if (this.showActions) {
      this.defaultHeaders.push(
        {
          text: 'Actions',
          value: 'action',
          sortable: false
        })
    }
    return this.defaultHeaders
  }

  @Prop({
    required: true,
    default: () => {
      return []
    }
  })
  accounts!: Account[]
  searched: Account[] = []
  search: string = ''
  itemsPerPage: number = 25

  @Prop({
    type: Boolean,
    default: false
  })
  loading!: boolean

  @Prop({
    type: Boolean,
    default: false
  })
  showGroup!: boolean

  @Prop({
    type: Boolean,
    default: false
  })
  showActions!: boolean

  @Prop({
    type: Boolean,
    default: false
  })
  showSearch!: boolean

  @Prop({
    type: Boolean,
    default: false
  })
  showTitle!: boolean

  @Watch('accounts')
  updateSearched() {
    this.searched = this.accounts
  }

  editItem(item: Account) {
    this.$router.push({ name: 'account', params: { id: item.id.toString() } })
  }

  deleteItem(item: Account) {
    const index = this.searched.indexOf(item)
    confirm('Are you sure you want to delete this item?') && this.searched.splice(index, 1)
  }

  created() {
    this.searched = this.accounts
  }
}
</script>

<style scoped lang="scss">
  tr {
    cursor: pointer;
  }
</style>
