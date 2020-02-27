<template>
  <div class="account-details">
    <v-card class="account-form">
      <v-card-title>Account {{account.name}}</v-card-title>
      <v-card-text>
        <account-form :account="account" v-if="account"/>
      </v-card-text>
    </v-card>
    <transaction-list :transactions="transactions"
                      :empty-message="`No transactions found for '${account.name}', maybe this account has none`" v-if="account"/>
  </div>
</template>

<script lang="ts">
import { Component, Vue } from 'vue-property-decorator'
import Account, { emptyAccount } from '@/data/account'
import TransactionList from '@/components/transactions/List.vue'
import Transactions from '@/views/Transactions.vue'
import Transaction from '@/data/transaction'
import AccountForm from '@/components/account/Form.vue'
import axios from 'axios'

@Component({
  components: { Transactions, TransactionList, AccountForm }
})
export default class AccountDetail extends Vue {
  account: Account = emptyAccount
  transactions: Transaction[] = []

  created() {
    axios.get('/account/' + this.$route.params.id, {
      baseURL: 'http://localhost:8088/v1',
      headers: {
        // @ts-ignore
        Authorization: 'Bearer ' + JSON.parse(localStorage.getItem('jwt')).access_token
      }
    }).then((response) => {
      let resAccount = response.data
      this.account = response.data
      if (!resAccount.hasOwnProperty('description')) {
        this.account.description = ''
      }
      this.account.nfcChipId = response.data.nfc_chip_id
    }).catch((response) => {
      console.error(response)
    })
    axios.get('/account/' + this.$route.params.id + '/transactions', {
      baseURL: 'http://localhost:8088/v1',
      headers: {
        // @ts-ignore
        Authorization: 'Bearer ' + JSON.parse(localStorage.getItem('jwt')).access_token
      }
    }).then((response) => {
      let transactions = response.data.transactions
      this.transactions = transactions.map((el) => {
        return {
          id: el.id,
          newSaldo: el.new_saldo,
          oldSaldo: el.old_saldo,
          account: el.account,
          created: el.created,
          amount: el.amount
        }
      })
    }).catch((response) => {
      console.error(response)
    })
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
