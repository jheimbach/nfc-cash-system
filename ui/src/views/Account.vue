<template>
  <div class="account-details">
    <v-card class="account-form">
      <v-card-title>Account {{account.name}}</v-card-title>
      <v-card-text>
        <account-form :account="account"/>
      </v-card-text>
    </v-card>
    <transaction-list :transactions="transactions"
                      :empty-message="`No transactions found for '${account.name}', maybe this account has none`"/>
  </div>
</template>

<script lang="ts">
import { Component, Vue } from 'vue-property-decorator'
import Account from '@/data/account'
import TransactionList from '@/components/transactions/List.vue'
import Transactions from '@/views/Transactions.vue'
import Transaction from '@/data/transaction'
import AccountForm from '@/components/account/Form.vue'

@Component({
  components: { Transactions, TransactionList, AccountForm }
})
export default class AccountDetail extends Vue {
  account!: Account
  transactions: Transaction[] = []

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
    this.transactions = [
      {
        id: 1,
        oldSaldo: 540,
        newSaldo: 539,
        amount: 1,
        created: new Date('2018-12-10T01:58:06Z'),
        account: this.account
      }, {
        id: 2,
        oldSaldo: 539,
        newSaldo: 532,
        amount: 7,
        created: new Date('2018-12-10T04:10:15Z'),
        account: this.account
      }
    ]
  }
}
</script>

<style scoped>
  .account-form {
    margin-bottom: 20px;
  }

  .account-details {
    width: 100%;
  }
</style>
