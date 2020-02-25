<template>
  <account-form :account="singleAccount" clear-btn="Reset Form" save-btn="Create" @save="save">
    <template v-slot:snackbar="{activator, hideSnackbar}">
      <v-snackbar :value="activator" @input="hideSnackbar">
        Account {{singleAccount.name }} successfully created
        <v-btn color="pink" text @click="hideSnackbar">Close</v-btn>
      </v-snackbar>
    </template>
  </account-form>
</template>

<script lang="ts">

import { Component, Vue } from 'vue-property-decorator'
import AccountForm from '@/components/account/Form.vue'
import Account, { emptyAccount } from '@/data/account'

@Component({
  components: { AccountForm }
})
export default class AccountCreateSingle extends Vue {
  singleAccount: Account = Object.assign({}, emptyAccount)

  save() {
    this.singleAccount = emptyAccount
  }
}
</script>

<style scoped>

</style>
