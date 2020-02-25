<template>
  <v-form ref="form" v-model="valid" lazy-validation>
    <v-text-field
      v-model="account.id"
      label="ID"
      disabled
      prepend-icon="vpn_key"
      v-if="account.id !== 0"
    />
    <v-text-field
      v-model="account.name"
      :counter="255"
      :rules="[(v) => !!v || 'Name is required']"
      label="Name"
      required
      prepend-icon="person"
    />
    <v-textarea
      v-model="account.description"
      label="Description"
      prepend-icon="subject"
    />
    <v-text-field
      v-model="account.saldo"
      type="number"
      label="Saldo"
      prepend-icon="attach_money"
      :disabled="account.id !== 0"
    />
    <v-select
      v-model="account.group"
      :items="groups"
      item-text="name"
      item-value="id"
      return-object
      :rules="[v => !!v || 'Group is required']"
      label="Group"
      required
      prepend-icon="group"
    />
    <v-text-field
      v-model="account.nfcChipId"
      :counter="20"
      :rules="chipRules"
      label="NFC Chip ID"
      required
      prepend-icon="nfc"
    />
    <v-btn :disabled="!valid" color="success" class="mr-4" @click="save">Save</v-btn>
    <v-btn color="error" class="mr-4" @click="cancel">Clear Changes</v-btn>
    <v-snackbar v-model="snackbar">
      Account {{account.name }} successfully saved
      <v-btn color="pink" text @click="snackbar = false">Close</v-btn>
    </v-snackbar>
  </v-form>
</template>

<script lang="ts">
import { Component, Emit, Prop, Vue } from 'vue-property-decorator'
import Account from '@/data/account'
import Group from '@/data/group'

@Component
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
  groups: Group[] = []
  unchangedAccount!: Account
  valid: boolean = true
  snackbar: boolean = false
  chipRules = [
    (v: string) => !!v || 'NFC Chip ID is required',
    (v: string) => (v || '').length <= 20 || 'NFC Chip ID should be max. 20 characters'
  ]

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

  created() {
    this.groups = [
      {
        id: 1,
        name: 'H2O Plus'
      },
      {
        id: 2,
        name: 'A-S Medication Solutions LLC',
        description: 'E.E.S',
        canOverdraw: true
      },
      {
        id: 3,
        name: 'Mylan Pharmaceuticals Inc.'
      },
      {
        id: 4,
        name: 'Mylan Pharmaceuticals Inc.',
        description: 'Enalapril Maleate and Hydrochlorothiazide'
      },
      {
        id: 5,
        name: 'REMEDYREPACK INC.',
        description: 'CELEBREX'
      },
      {
        id: 6,
        name: 'H E B',
        description: 'night time'
      },
      {
        id: 7,
        name: 'PSS World Medical, Inc.',
        canOverdraw: true
      },
      {
        id: 8,
        name: 'Kareway Product, Inc.',
        description: 'Acetaminophen',
        canOverdraw: true
      },
      {
        id: 9,
        name: 'Pharmacia and Upjohn Company'
      },
      {
        id: 10,
        name: 'Dolgencorp, Inc. (DOLLAR GENERAL & REXALL)',
        description: 'Allergy Relief',
        canOverdraw: true
      }
    ]
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
