<template>
  <v-card class="account-create">
    <v-card-title>Create new Account</v-card-title>
    <v-card-text v-if="creationType('single')">
      <account-form :account="singleAccount" clear-btn="Reset Form" save-btn="Create" @save="save">
        <template v-slot:snackbar="{activator, hideSnackbar}">
          <v-snackbar :value="activator" @input="hideSnackbar">
            Account {{singleAccount.name }} successfully created
            <v-btn color="pink" text @click="hideSnackbar">Close</v-btn>
          </v-snackbar>
        </template>
      </account-form>
    </v-card-text>
    <v-card-text v-if="creationType('multiple')">
      <account-form-multiple :accounts="accounts" />
    </v-card-text>
    <v-card-text v-if="creationType('upload')">
      <v-file-input label="Upload CSV-File" show-size prepend-icon="file_copy"></v-file-input>
    </v-card-text>
  </v-card>
</template>

<script lang="ts">

import { Component, Vue, Watch } from 'vue-property-decorator'
import AccountForm from '@/components/account/Form.vue'
import Account, { emptyAccount } from '@/data/account'
import AccountFieldName from '@/components/form/Account/Name.vue'
import AccountFieldDescription from '@/components/form/Account/Description.vue'
import AccountFieldNfcChipId from '@/components/form/Account/NfcChipId.vue'
import AccountFieldGroup from '@/components/form/Account/Group.vue'
import AccountFieldSaldo from '@/components/form/Account/Saldo.vue'
import AccountFormMultiple from '@/components/account/FormMultiple.vue'

@Component({
  components: { AccountForm, AccountFormMultiple }
})
export default class AccountEdit extends Vue {
  accounts: Account[] = []
  singleAccount: Account = Object.assign({}, emptyAccount)

  save() {
    this.singleAccount = emptyAccount
  }

  creationType(typename: string) {
    if (typename === 'multiple') {
      return this.$route.query.multiple
    }
    if (typename === 'upload') {
      return this.$route.query.upload
    }
    if (typename === 'single') {
      return !this.$route.query.hasOwnProperty('multiple') && !this.$route.query.hasOwnProperty('upload')
    }
    return false
  }
}
</script>

<style scoped lang="scss">
  .account-create {
    width: 100%;
  }
</style>
