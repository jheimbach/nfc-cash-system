<template>
  <v-container>
    <v-row>
      <v-col cols="12" md="12">
        <v-text-field
          v-model="accountLength"
          type="number"
          label="No. of Accounts"
          required
          min="0"
          hint="how many accounts do you want to be created?"
        ></v-text-field>
      </v-col>
    </v-row>
    <v-form v-model="valid[index]" ref="form" class="form" v-for="(acc, index) in models" :key="'account'+index">
      <v-row align="center">
        <v-col cols="12" md="2">
          <account-field-name v-model="acc.name"/>
        </v-col>

        <v-col cols="12" md="3">
          <account-field-description v-model="acc.description"/>
        </v-col>

        <v-col cols="12" md="2">
          <account-field-saldo v-model="acc.saldo"/>
        </v-col>

        <v-col cols="12" md="2">
          <account-field-group v-model="acc.group"/>
        </v-col>

        <v-col cols="12" md="2">
          <account-field-nfc-chip-id v-model="acc.nfcChipId"/>
        </v-col>
        <v-col>
          <v-btn icon @click="removeAccount(index)">
            <v-icon>delete</v-icon>
          </v-btn>
        </v-col>
      </v-row>
    </v-form>
    <v-row>
      <v-col cols="12" md="12">
        <v-btn :disabled="!isValid()" color="success" class="mr-4" @click="save">Create Accounts</v-btn>
      </v-col>
    </v-row>
    <slot name="snackbar" v-bind:activator="snackbar" v-bind:hideSnackbar="hideSnackbar">
      <v-snackbar :value="snackbar" @input="hideSnackbar">
        Accounts successfully created
        <v-btn color="pink" text @click="hideSnackbar">Close</v-btn>
      </v-snackbar>
    </slot>
  </v-container>
</template>

<script lang="ts">
import { Component, Prop, Vue, Watch, Emit } from 'vue-property-decorator'
import AccountFieldName from '../form/Account/Name.vue'
import AccountFieldDescription from '../form/Account/Description.vue'
import AccountFieldNfcChipId from '../form/Account/NfcChipId.vue'
import AccountFieldGroup from '../form/Account/Group.vue'
import AccountFieldSaldo from '../form/Account/Saldo.vue'
import Account, { emptyAccount } from '@/data/account'

@Component({
  components: { AccountFieldName, AccountFieldDescription, AccountFieldNfcChipId, AccountFieldGroup, AccountFieldSaldo }
})
export default class AccountFormMultiple extends Vue {
  @Prop({
    type: Array,
    required: true
  })
  accounts!: Account[]
  models: Account[] = []
  valid: boolean[] = []
  accountLength: number = 5
  snackbar: boolean = false

  @Watch('accountLength')
  updateAccountLength() {
    const numberOfAccounts = this.accountLength - this.models.length
    if (numberOfAccounts > 0) {
      for (let i = 0; i < numberOfAccounts; i++) {
        this.models.push(Object.assign({}, emptyAccount))
        this.valid.push(false)
      }
    }
    if (numberOfAccounts < 0) {
      this.models = this.models.slice(0, numberOfAccounts)
      this.valid = this.valid.slice(0, numberOfAccounts)
    }
  }

  removeAccount(idx: number) {
    this.models.splice(idx, 1)
    this.valid.splice(idx, 1)
    this.accountLength--
  }

  created() {
    this.updateAccountLength()
  }

  hideSnackbar() {
    this.snackbar = false
  }

  isValid() {
    for (let i = 0; i < this.valid.length; i++) {
      if (!this.valid[i]) {
        return false
      }
    }
    return true
  }

  save() {
    // @ts-ignore
    const isValid = this.$refs.form.reduce((all, el) => {
      if (!el.validate()) {
        all = false
      }
      return all
    }, true)
    if (isValid) {
      this.snackbar = true
      this.$emit('save', this.models.filter((el) => {
        return JSON.stringify(el) !== JSON.stringify(emptyAccount)
      }))
    }
  }
}
</script>

<style scoped lang="scss">
  .form {
    width: 100%;
  }
</style>
