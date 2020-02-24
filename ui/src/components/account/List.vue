<template>
  <v-card width="100%">
    <v-card-title>Accounts
      <v-spacer/>
      <v-text-field v-model="search" append-icon="search" label="Search" single-line hide-details/>
    </v-card-title>
    <v-data-table :headers="headers" :items="searched" :items-per-page="50" fixed-header item-key="id" :search="search"
                  :loading="loading">
      <template v-slot:item.action="{ item }">
        <v-icon
          small
          class="mr-2"
          @click="editItem(item)"
        >
          edit
        </v-icon>
        <v-icon
          small
          @click="deleteItem(item)"
        >
          delete
        </v-icon>
      </template>
    </v-data-table>
  </v-card>
</template>

<script lang="ts">

import { Component, Emit, Prop, Vue } from 'vue-property-decorator'
import Account from '@/data/account'
import { DataTableHeader } from 'vuetify'

@Component
export default class AccountList extends Vue {
  headers: DataTableHeader[] = [
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
      text: 'Group',
      value: 'group.name',
      sortable: false
    },
    {
      text: 'NFC Chip ID',
      align: 'end',
      value: 'nfcChipId',
      sortable: false
    },
    {
      text: 'Actions',
      value: 'action',
      sortable: false
    }
  ]
  @Prop({
    required: true,
    default: () => {
      return []
    }
  })
  accounts!: Account[]
  searched: Account[] = []
  search: string = ''

  @Prop({
    type: Boolean,
    default: false
  })
  loading!: boolean

  editItem(item: Account) {
    this.$router.push({ name: 'account', params: { id: item.id.toString() } })
  }

  deleteItem(item: Account) {
    this.searched = this.accounts.filter((el) => {
      return el.id !== item.id
    })
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
