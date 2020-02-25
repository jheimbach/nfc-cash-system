<template>
  <div class="group-details">
    <v-card class="group-form">
      <v-card-title>Group {{group.name}}</v-card-title>
      <v-card-text>
        <group-form :group="group"/>
      </v-card-text>
    </v-card>
    <account-list :accounts="accounts" show-actions show-title/>
  </div>
</template>

<script lang="ts">
import { Component, Vue } from 'vue-property-decorator'
import AccountList from '@/components/account/List.vue'
import GroupForm from '@/components/groups/Form.vue'
import Group from '@/data/group'
import Account from '@/data/account'

@Component({
  components: { AccountList, GroupForm }
})
export default class GroupDetail extends Vue {
  group: Group = { id: 0, name: '' }
  accounts: Account[] = []

  created() {
    this.group = {
      id: 1,
      name: 'H20 Plus',
      description: 'test',
      canOverdraw: false
    }
    this.accounts = [
      {
        id: 13,
        name: 'Kristin Cicullo',
        description: 'B.S.C AMPUL',
        saldo: 373,
        nfcChipId: 'DseWgH8AA1l',
        group: {
          id: 1,
          name: 'H2O Plus'
        }
      },
      {
        id: 18,
        name: 'Aluin Cunnah',
        description: 'KADIAN',
        saldo: 361,
        nfcChipId: 'WNPZMY',
        group: {
          id: 1,
          name: 'H2O Plus'
        }
      },
      {
        id: 19,
        name: 'Gabriell Nunnerley',
        saldo: 431,
        nfcChipId: 'Jzdvpr',
        group: {
          id: 1,
          name: 'H2O Plus'
        }
      },
      {
        id: 32,
        name: 'Renaud Delacroux',
        saldo: 496,
        nfcChipId: 'nieCIfe',
        group: {
          id: 1,
          name: 'H2O Plus'
        }
      },
      {
        id: 35,
        name: 'Geoff Cornejo',
        saldo: 439,
        nfcChipId: 'Uo57KDmf',
        group: {
          id: 1,
          name: 'H2O Plus'
        }
      },
      {
        id: 37,
        name: 'Hugh Cunnell',
        description: 'Nitrotan',
        saldo: 466,
        nfcChipId: 'eJZHaCH',
        group: {
          id: 1,
          name: 'H2O Plus'
        }
      },
      {
        id: 65,
        name: 'Roana Grady',
        description: 'Thermazene',
        saldo: 442,
        nfcChipId: 'u1egQQtK',
        group: {
          id: 1,
          name: 'H2O Plus'
        }
      },
      {
        id: 73,
        name: 'Merlina Gallagher',
        saldo: 485,
        nfcChipId: 'y0uh7p5lrc',
        group: {
          id: 1,
          name: 'H2O Plus'
        }
      },
      {
        id: 77,
        name: 'Randa Bernard',
        description: 'Corn Smut',
        saldo: 484,
        nfcChipId: 'rShbgKX3',
        group: {
          id: 1,
          name: 'H2O Plus'
        }
      },
      {
        id: 79,
        name: 'Kaela Crosser',
        saldo: 458,
        nfcChipId: 'PtjPRh',
        group: {
          id: 1,
          name: 'H2O Plus'
        }
      }
    ]
  }
}

</script>

<style scoped lang="scss">
  .group-details {
    width: 100%;
  }

  .group-form {
    width: 100%;
    margin-bottom: 20px;
  }
</style>
