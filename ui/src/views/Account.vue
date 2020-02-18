<template>
  <div>
    <h1>Account {{account.name}}</h1>
    <md-tabs>
      <md-tab id="tab-details" md-label="Account details">
        <account-detail-list :account="account" :edit="editAccount"/>
      </md-tab>
      <md-tab id="tab-transactions" md-label="Transactions">
        <h2>Transactionlist</h2>
        <transaction-list/>
        <!--<transaction-list :transactions="transactions"/>-->
      </md-tab>
    </md-tabs>
  </div>
</template>

<script lang="ts">
import { Component, Vue } from 'vue-property-decorator'
import Account from '@/data/account'
import AccountDetailList from '@/components/account/Details.vue'
import TransactionList from '@/components/transactions/list.vue'
import Transactions from '@/views/Transactions.vue'

@Component({
  components: { Transactions, AccountDetailList, TransactionList }
})
export default class AccountCmp extends Vue {
  editAccount: boolean = false
  account!: Account

  created() {
    this.account = {
      id: 1,
      name: 'test',
      description: 'Lorem ipsum dolor sit amet, Lorem ipsum dolor sit amet, consectetur adipisicing elit. Lorem ipsum dolor sit amet, consectetur adipisicing elit. Lorem ipsum dolor sit amet, consectetur adipisicing elit. Lorem ipsum dolor sit amet, consectetur adipisicing elit.',
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
