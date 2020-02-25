<template>
  <v-card class="account-create">
    <v-card-title>Create new Account</v-card-title>
    <v-card-text>
      <router-view/>
    </v-card-text>
  </v-card>
</template>

<script lang="ts">

import { Component, Vue } from 'vue-property-decorator'

@Component
export default class AccountCreate extends Vue {
}
</script>

<style scoped lang="scss">
  .account-create {
    width: 100%;
  }
</style>
