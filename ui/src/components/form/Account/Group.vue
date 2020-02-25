<template>
  <v-select
    v-bind:value="value"
    @input="$emit('input', $event)"
    :items="groups"
    item-text="name"
    item-value="id"
    return-object
    :rules="[v => !!v || 'Group is required']"
    label="Group"
    required
    prepend-icon="group"
  />
</template>

<script lang="ts">
import { Component, Prop, Vue } from 'vue-property-decorator'
import Group from '@/data/group'

@Component
export default class AccountFieldGroup extends Vue {
  @Prop({
    required: true
  })
  value!: Group | null
  groups: Group[] = []

  created() {
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
}
</script>

<style scoped>

</style>
