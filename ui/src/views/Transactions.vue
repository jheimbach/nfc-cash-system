<template>
  <transaction-list :transactions="transactions" show-actions show-search show-account-name :loading="loading"/>
</template>

<script lang="ts">
import { Component, Vue } from 'vue-property-decorator'
import TransactionList from '@/components/transactions/List.vue'
import Transaction from '@/data/transaction'
import axios from 'axios'

@Component({
  components: { TransactionList }
})
export default class TransactionsView extends Vue {
  transactions: Transaction[] = []
  loading: boolean = false

  getTransactions() {
    axios.get('/transactions', {
      baseURL: 'http://localhost:8088/v1',
      headers: {
        Authorization: 'Bearer ' + JSON.parse(localStorage.getItem('jwt')).access_token
      }
    }).then((response) => {
      this.transactions = response.data.transactions.map((el) => {
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

  created() {
    this.loading = true
    setTimeout(() => {
      this.loading = false
      this.getTransactions()
    }, 1000)
  }
}
</script>

<style lang="scss" scoped>
</style>
