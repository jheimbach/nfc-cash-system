<template>
  <account-form-multiple :accounts="accounts" />
</template>

<script lang="ts">

import { Component, Vue } from 'vue-property-decorator'
import Account from '@/data/account'
import AccountFormMultiple from '@/components/account/FormMultiple.vue'

@Component({
  components: { AccountFormMultiple }
})
export default class AccountCreateMultiple extends Vue {
  accounts: Account[] = []
}
</script>
<style scoped>

</style>
