<template>
  <div class="account-form">
    <div class="row">
      <md-field>
        <label>Name</label>
        <md-input v-model="model.name"></md-input>
      </md-field>
    </div>
    <div class="row">
      <md-field>
        <label>Description</label>
        <md-textarea v-model="model.description"></md-textarea>
      </md-field>
    </div>
    <div class="row" v-if="model.id === 0">
      <md-field>
        <label>Saldo</label>
        <md-input v-model="model.saldo" type="number"></md-input>
      </md-field>
    </div>
    <div class="row">
      <md-field>
        <label>NFC Chip ID</label>
        <md-input v-model="model.nfcChipId"></md-input>
      </md-field>
    </div>
    <div class="row">
      <md-field>
        <label>Group</label>
        <md-select v-model="groupModel" name="group" id="group" @md-selected="groupIdToGroup">
          <md-option :value="group.id" v-for="group in groups" v-bind:key="group.id">{{group.name}}</md-option>
        </md-select>
      </md-field>
    </div>
    <div class="row button-row" v-if="!noButtons">
      <md-button @click="cancel" class="md-raised md-accent">Cancel</md-button>
      <md-button @click="save" class="md-raised md-primary">Save</md-button>
    </div>
  </div>
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
  @Prop({
    default: false,
    type: Boolean
  })
  noButtons!: boolean
  model!:Account
  groups: Group[] = []
  mappedGroups: any = {}
  groupModel: number = 0
  unchangedAccount!: Account

  @Emit('save')
  save(): Account {
    return JSON.parse(JSON.stringify(this.model))
  }

  @Emit('cancel')
  cancel() {
    const js = JSON.parse(JSON.stringify(this.unchangedAccount))
    this.model = js
    return js
  }

  groupIdToGroup(val: number) {
    this.model.group = this.mappedGroups[val]
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
    this.groups.forEach((el) => {
      this.mappedGroups[el.id] = el
    })
    this.model = this.account
    this.unchangedAccount = JSON.parse(JSON.stringify(this.model))
    this.groupModel = this.model.group.id
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
