<template>
  <div>
    <md-toolbar md-elevation="1">
      <h3 class="md-title" style="flex: 1">Account {{account.name}}</h3>
      <md-button class="md-primary" @click="openEditView">Edit</md-button>
    </md-toolbar>
    <account-detail-list :account="account"/>
    <transaction-list :transactions="transactions"
                      :empty-message="`No transactions found for '${account.name}', maybe this account has none`"/>

    <md-dialog :md-active.sync="editView">
      <md-dialog-title>Edit account {{account.name}}</md-dialog-title>
      <account-form :account="account" @cancel="closeEditView" @save="closeEditView" ref="form" no-buttons/>
      <md-dialog-actions>
        <md-button class="md-primary" @click="passCancelToForm()">Close</md-button>
        <md-button class="md-primary" @click="passSaveToForm()">Save</md-button>
      </md-dialog-actions>
    </md-dialog>
  </div>
</template>

<script lang="ts">
import { Component, Vue } from 'vue-property-decorator'
import Account from '@/data/account'
import AccountDetailList from '@/components/account/Details.vue'
import TransactionList from '@/components/transactions/List.vue'
import Transactions from '@/views/Transactions.vue'
import { Transaction } from '@/data/transaction'
import AccountForm from '@/components/account/Form.vue'

@Component({
  components: { Transactions, AccountDetailList, TransactionList, AccountForm }
})
export default class AccountDetail extends Vue {
  account!: Account
  transactions: Transaction[] = []
  editView: boolean = false

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

  passSaveToForm() {
    const form = this.$refs['form'] as AccountForm
    form.save()
  }

  passCancelToForm() {
    const form = this.$refs['form'] as AccountForm
    form.cancel()
  }

  closeEditView(account: Account) {
    this.account = account
    this.editView = false
  }

  openEditView() {
    this.editView = true
  }
}
</script>

<style scoped>
  .account-form {
    width: 75vw;
    padding: 20px;
  }
</style>
