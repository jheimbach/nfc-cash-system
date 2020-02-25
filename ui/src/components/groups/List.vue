<template>
  <v-card width="100%">
    <v-card-title>
      <template v-if="showTitle">Groups</template>
      <v-spacer v-if="showTitle && showSearch"/>
      <v-text-field v-model="search" append-icon="search" label="Search" single-line hide-details v-if="showSearch"/>
    </v-card-title>
    <v-data-table :headers="headers" :items="searched" :items-per-page="itemsPerPage" fixed-header item-key="id"
                  :search="search"
                  :loading="loading" :hide-default-footer="searched.length < itemsPerPage">
      <template v-slot:item.action="{ item }">
        <v-icon small class="mr-2" @click="editItem(item)">edit</v-icon>
        <v-icon small @click="deleteItem(item)">delete</v-icon>
      </template>
      <template v-slot:item.canOverdraw="{ item }">
        {{item.canOverdraw ? 'Yes':'No'}}
      </template>
    </v-data-table>
  </v-card>
</template>

<script lang="ts">
import { Component, Prop, Vue, Watch } from 'vue-property-decorator'
import Group from '@/data/group'
import { DataTableHeader } from 'vuetify'

@Component
export default class GroupList extends Vue {
  @Prop({
    default: () => {
      return []
    },
    required: true,
    type: Array
  })
  groups!: Group[]
  @Prop({
    default: false,
    type: Boolean
  })
  loading!: boolean
  @Prop({
    default: false,
    type: Boolean
  })
  showTitle!: boolean
  @Prop({
    default: false,
    type: Boolean
  })
  showSearch!: boolean
  searched: Group[] = []
  search: string = ''
  itemsPerPage: number = 25

  headers: DataTableHeader[] = [
    {
      text: 'ID',
      value: 'id'
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
      text: 'Can Overdraw',
      value: 'canOverdraw'
    },
    {
      text: 'Actions',
      value: 'action',
      sortable: false
    }
  ]

  @Watch('groups')
  updateSearched() {
    this.searched = this.groups
  }

  editItem(item: Group) {
    this.$router.push({ name: 'group', params: { id: item.id.toString() } })
  }

  deleteItem(item: Group) {
    const index = this.searched.indexOf(item)
    confirm(`Are you sure you want to delete ${item.name}?`) && this.searched.splice(index, 1)
  }

  created() {
    this.searched = this.groups
  }
}
</script>

<style scoped lang="scss">

</style>
