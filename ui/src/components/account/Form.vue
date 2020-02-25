<template>
  <v-form ref="form" v-model="valid" lazy-validation>
    <v-text-field
      v-model="account.id"
      label="ID"
      disabled
      prepend-icon="vpn_key"
      v-if="account.id !== 0"
    />
    <account-field-name v-model="account.name"/>
    <account-field-description v-model="account.description" textarea/>
    <account-field-saldo v-model="account.saldo" :disabled="account.id !== 0"/>
    <account-field-group v-model="account.group"/>
    <account-field-nfc-chip-id v-model="account.nfcChipId"/>
    <v-btn :disabled="!valid" color="success" class="mr-4" @click="save">{{saveBtn}}</v-btn>
    <v-btn color="error" class="mr-4" @click="cancel">{{clearBtn}}</v-btn>
    <slot name="snackbar" v-bind:activator="snackbar" v-bind:hideSnackbar="hideSnackbar">
      <v-snackbar :value="snackbar" @input="hideSnackbar">
        Account {{account.name }} successfully saved
        <v-btn color="pink" text @click="hideSnackbar">Close</v-btn>
      </v-snackbar>
    </slot>
  </v-form>
</template>

<script lang="ts">
import { Component, Emit, Prop, Vue } from 'vue-property-decorator'
import Account from '@/data/account'
import AccountFieldName from '@/components/form/Account/Name.vue'
import AccountFieldDescription from '@/components/form/Account/Description.vue'
import AccountFieldNfcChipId from '@/components/form/Account/NfcChipId.vue'
import AccountFieldGroup from '@/components/form/Account/Group.vue'
import AccountFieldSaldo from '@/components/form/Account/Saldo.vue'

@Component({
  components: { AccountFieldName, AccountFieldDescription, AccountFieldNfcChipId, AccountFieldGroup, AccountFieldSaldo }
})
export default class AccountForm extends Vue {
  @Prop({
    default: () => {
      return {
        id: 0,
        name: '',
        description: '',
        saldo: 0,
        nfcChipId: '',
        group: null
      }
    },
    required: true,
    type: Object
  })
  account!: Account
  @Prop({
    type: String,
    default: 'Clear Changes'
  })
  clearBtn!: string
  @Prop({
    type: String,
    default: 'Save'
  })
  saveBtn!: string
  unchangedAccount!: Account
  valid: boolean = true
  snackbar: boolean = false

  @Emit('save')
  save() {
    // @ts-ignore
    if (this.$refs.form.validate()) {
      setTimeout(() => {
        this.snackbar = true
      }, 500)
    }
    return this.account
  }

  @Emit('cancel')
  cancel() {
    this.account.name = this.unchangedAccount.name
    this.account.description = this.unchangedAccount.description
    this.account.nfcChipId = this.unchangedAccount.nfcChipId
    this.account.group = this.unchangedAccount.group
    return this.account
  }

  hideSnackbar() {
    this.snackbar = false
  }

  created() {
    this.unchangedAccount = JSON.parse(JSON.stringify(this.account))
  }
}
</script>

<style scoped lang="scss">
  .row {
    display: flex;

    &.button-row {
      flex-direction: row;
      justify-content: flex-end;
    }
  }
</style>
