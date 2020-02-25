<template>
  <multi-model-form :models="accounts" @lenUpdated="updateAccounts" @removeModel="removeAccount"/>
</template>

<script lang="ts">

import { Component, Vue } from 'vue-property-decorator'
import Account, { emptyAccount } from '@/data/account'
import MultiModelForm from '@/components/MultiModelForm.vue'

@Component({
  components: { MultiModelForm }
})
export default class AccountCreateMultiple extends Vue {
  accounts: Account[] = []

  updateAccounts(len: number) {
    const numberOfAccounts = len - this.accounts.length
    if (numberOfAccounts > 0) {
      for (let i = 0; i < numberOfAccounts; i++) {
        this.accounts.push(Object.assign({}, emptyAccount))
      }
    }
    if (numberOfAccounts < 0) {
      this.accounts = this.accounts.slice(0, numberOfAccounts)
    }
  }

  removeAccount(idx: number) {
    this.accounts.splice(idx, 1)
  }
}
</script>
<style scoped>

</style>
