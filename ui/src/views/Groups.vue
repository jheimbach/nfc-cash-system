<template>
  <group-list :groups="groups" :loading="loading"/>
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
</style>
