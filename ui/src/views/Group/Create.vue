<template>
  <v-card class="group-create">
    <v-card-title>Create new Group</v-card-title>
    <v-card-text>
      <router-view/>
    </v-card-text>
  </v-card>
</template>

<script lang="ts">

import { Component, Vue } from 'vue-property-decorator'

@Component
export default class GroupCreate extends Vue {
}
</script>

<style scoped lang="scss">
  .group-create {
    width: 100%;
  }
</style>
