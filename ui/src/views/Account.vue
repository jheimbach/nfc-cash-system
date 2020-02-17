<template>
  <div>
    <account-detail-list :account="account"/>
    <!--<transaction-list :transactions="transactions"/>-->
  </div>
</template>

<script lang="ts">
import { Component, Vue } from 'vue-property-decorator'
import Account from '@/data/account'
import AccountDetailList from '@/components/AccountDetailList.vue'

@Component({
  components: { AccountDetailList }
})
export default class AccountCmp extends Vue {
  account: Account = null

  created() {
    this.account = {
      id: 1,
      name: 'test',
      description: 'test desc',
      saldo: 12,
      nfcChipId: 'asdfwd',
      group: {
        id: 1,
        name: 'group one'
      }
    }
  }
}
</script>

<style scoped lang="scss">
</style>
