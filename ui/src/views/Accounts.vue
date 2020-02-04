<template>
  <div>
    <md-table v-model="searched" md-sort="name" md-sort-order="asc">
      <md-table-toolbar>
        <div class="md-toolbar-section-start">
          <h1 class="md-title">Accounts</h1>
        </div>

        <md-field md-clearable class="md-toolbar-section-end">
          <md-input placeholder="Search by name..." v-model="search" @input="searchOnTable"/>
        </md-field>
      </md-table-toolbar>

      <md-table-empty-state
        md-label="No accounts found"
        :md-description="`No accounts found for this '${search}' query. Try a different search term.`">
      </md-table-empty-state>

      <md-table-row slot="md-table-row" slot-scope="{ item }">
        <md-table-cell md-label="ID" md-sort-by="id" md-numeric>{{ item.id }}</md-table-cell>
        <md-table-cell md-label="Name" md-sort-by="name">{{ item.name }}</md-table-cell>
        <md-table-cell md-label="Description" md-sort-by="email">{{ item.description }}</md-table-cell>
        <md-table-cell md-label="Saldo" md-sort-by="gender">{{ item.saldo }}</md-table-cell>
        <md-table-cell md-label="Group" md-sort-by="title">{{ item.group.name }}</md-table-cell>
        <md-table-cell md-label="NFC Chip Id" md-sort-by="title">{{ item.nfcChipId }}</md-table-cell>
        <md-table-cell md-label="Actions">
          <md-button :to="transactionsRoute(item.id)" class="md-icon-button ncs-secondary">
            <md-icon>account_balance_wallet</md-icon>
          </md-button>
          <md-button :to="editRoute(item.id)" class="md-icon-button md-primary">
            <md-icon>edit</md-icon>
          </md-button>
          <delete-confirmation :title="`Delete Account ${item.name}?`"
                               content="Do you really want to delete this account? This action cannot be reversed"
                               @onDelete="deleteAccount(item)"
          >

          </delete-confirmation>
        </md-table-cell>
      </md-table-row>
    </md-table>
  </div>
</template>

<script lang="ts">

import { Component, Prop, Vue } from 'vue-property-decorator'
import DeleteConfirmation from '@/components/DeleteConfirmation.vue'
import Account from '@/data/account'

@Component({
  components: { DeleteConfirmation }
})

export default class Accounts extends Vue {
  accounts: Account[] = []
  search: string = ''
  searched: Account[] = []

  getAccounts() {
    this.accounts = [
      {
        id: 1,
        name: 'Laverne Blackstock',
        description: 'Itchy Eye',
        saldo: 436,
        nfcChipId: 'Hv8mnajqzIKO',
        group: {
          id: 7,
          name: 'PSS World Medical, Inc.',
          canOverdraw: true
        }
      },
      {
        id: 2,
        name: 'Misha Blowfelde',
        saldo: 449,
        nfcChipId: '0XPPQy4ZkO7',
        group: {
          id: 5,
          name: 'REMEDYREPACK INC.',
          description: 'CELEBREX'
        }
      },
      {
        id: 3,
        name: 'Winnie Rennolds',
        description: 'Ofloxacin',
        saldo: 436,
        nfcChipId: 'ofzGN0eS2K',
        group: {
          id: 6,
          name: 'H E B',
          description: 'night time'
        }
      },
      {
        id: 4,
        name: 'Gordy Johnson',
        saldo: 462,
        nfcChipId: 'KTehLLhT',
        group: {
          id: 3,
          name: 'Mylan Pharmaceuticals Inc.'
        }
      },
      {
        id: 5,
        name: 'Kessia Spadollini',
        description: 'Alcohol Prep Pad',
        saldo: 421,
        nfcChipId: 'YxN57MH6',
        group: {
          id: 8,
          name: 'Kareway Product, Inc.',
          description: 'Acetaminophen',
          canOverdraw: true
        }
      },
      {
        id: 6,
        name: 'Haley Waker',
        saldo: 487,
        nfcChipId: '5putPvT',
        group: {
          id: 7,
          name: 'PSS World Medical, Inc.',
          canOverdraw: true
        }
      },
      {
        id: 7,
        name: 'Yolanda Pelos',
        description: 'Rodan And Fields Essentials Lip Shield SPF 25',
        saldo: 523,
        nfcChipId: 'bFJxlWF',
        group: {
          id: 5,
          name: 'REMEDYREPACK INC.',
          description: 'CELEBREX'
        }
      },
      {
        id: 8,
        name: 'Melisa Josowitz',
        saldo: 461,
        nfcChipId: 'FsBhsEwr',
        group: {
          id: 9,
          name: 'Pharmacia and Upjohn Company'
        }
      },
      {
        id: 9,
        name: 'Matias Beininck',
        saldo: 466,
        nfcChipId: '0E9wTFbJ',
        group: {
          id: 8,
          name: 'Kareway Product, Inc.',
          description: 'Acetaminophen',
          canOverdraw: true
        }
      },
      {
        id: 10,
        name: 'Rachele Steptowe',
        saldo: 383,
        nfcChipId: 'D54UACeRRMNJ',
        group: {
          id: 8,
          name: 'Kareway Product, Inc.',
          description: 'Acetaminophen',
          canOverdraw: true
        }
      },
      {
        id: 11,
        name: 'Lynea Habberjam',
        saldo: 449,
        nfcChipId: 'uSe3Sj',
        group: {
          id: 9,
          name: 'Pharmacia and Upjohn Company'
        }
      },
      {
        id: 12,
        name: 'Jarret Herculson',
        saldo: 442,
        nfcChipId: 'V5SadofJww3',
        group: {
          id: 4,
          name: 'Mylan Pharmaceuticals Inc.',
          description: 'Enalapril Maleate and Hydrochlorothiazide'
        }
      },
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
        id: 14,
        name: 'Berta Radborne',
        saldo: 452,
        nfcChipId: 'oz2nuvTUMK',
        group: {
          id: 9,
          name: 'Pharmacia and Upjohn Company'
        }
      },
      {
        id: 15,
        name: 'Virgina Stairmond',
        saldo: 400,
        nfcChipId: 'nt9th6L7eD',
        group: {
          id: 9,
          name: 'Pharmacia and Upjohn Company'
        }
      },
      {
        id: 16,
        name: 'Rubi Howey',
        description: 'Perrigo Hydroquinone',
        saldo: 508,
        nfcChipId: 'aagVbL',
        group: {
          id: 4,
          name: 'Mylan Pharmaceuticals Inc.',
          description: 'Enalapril Maleate and Hydrochlorothiazide'
        }
      },
      {
        id: 17,
        name: 'Shanna Ace',
        description: 'Glipizide',
        saldo: 466,
        nfcChipId: 'RIL1IVl',
        group: {
          id: 9,
          name: 'Pharmacia and Upjohn Company'
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
        id: 20,
        name: 'Florida Duesberry',
        saldo: 495,
        nfcChipId: 'rKlqNQsxt',
        group: {
          id: 9,
          name: 'Pharmacia and Upjohn Company'
        }
      },
      {
        id: 21,
        name: 'Turner Tutton',
        saldo: 424,
        nfcChipId: 'LECCoD',
        group: {
          id: 2,
          name: 'A-S Medication Solutions LLC',
          description: 'E.E.S',
          canOverdraw: true
        }
      },
      {
        id: 22,
        name: 'Abbi Usher',
        description: 'Pravastatin Sodium',
        saldo: 485,
        nfcChipId: 'hXTuAvFk',
        group: {
          id: 8,
          name: 'Kareway Product, Inc.',
          description: 'Acetaminophen',
          canOverdraw: true
        }
      },
      {
        id: 23,
        name: 'Winona Pebworth',
        saldo: 440,
        nfcChipId: 'LouoW9Joku',
        group: {
          id: 9,
          name: 'Pharmacia and Upjohn Company'
        }
      },
      {
        id: 24,
        name: 'Horace Barnewell',
        description: 'Trout',
        saldo: 407,
        nfcChipId: 'tbaoWaYXbc',
        group: {
          id: 6,
          name: 'H E B',
          description: 'night time'
        }
      },
      {
        id: 25,
        name: 'Garrard Dreakin',
        saldo: 346,
        nfcChipId: 'I7oixdv',
        group: {
          id: 9,
          name: 'Pharmacia and Upjohn Company'
        }
      },
      {
        id: 26,
        name: 'Stafford Brewin',
        description: 'Constitutional Enhancer',
        saldo: 431,
        nfcChipId: 'dLz85N6',
        group: {
          id: 3,
          name: 'Mylan Pharmaceuticals Inc.'
        }
      },
      {
        id: 27,
        name: 'Kizzee Pinhorn',
        description: 'Cysto-Conray II',
        saldo: 439,
        nfcChipId: 'E3Fnu0oKgBC',
        group: {
          id: 4,
          name: 'Mylan Pharmaceuticals Inc.',
          description: 'Enalapril Maleate and Hydrochlorothiazide'
        }
      },
      {
        id: 28,
        name: 'Margaret Richie',
        description: 'Eye Irrigating',
        saldo: 479,
        nfcChipId: '0H34T9NR',
        group: {
          id: 2,
          name: 'A-S Medication Solutions LLC',
          description: 'E.E.S',
          canOverdraw: true
        }
      },
      {
        id: 29,
        name: 'Adey Ferfulle',
        description: 'Acetylcholine Chloride',
        saldo: 429,
        nfcChipId: 'PXJfRBm',
        group: {
          id: 2,
          name: 'A-S Medication Solutions LLC',
          description: 'E.E.S',
          canOverdraw: true
        }
      },
      {
        id: 30,
        name: 'Rollins Fullard',
        saldo: 424,
        nfcChipId: '2OKpbD3',
        group: {
          id: 10,
          name: 'Dolgencorp, Inc. (DOLLAR GENERAL & REXALL)',
          description: 'Allergy Relief',
          canOverdraw: true
        }
      },
      {
        id: 31,
        name: 'Wallache Bachelor',
        description: 'Clindamycin Hydrochloride',
        saldo: 430,
        nfcChipId: 'htiB3siU',
        group: {
          id: 9,
          name: 'Pharmacia and Upjohn Company'
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
        id: 33,
        name: 'Cassie Praundlin',
        saldo: 448,
        nfcChipId: 'Y8OSKZcIX9',
        group: {
          id: 9,
          name: 'Pharmacia and Upjohn Company'
        }
      },
      {
        id: 34,
        name: 'Alida Burkert',
        saldo: 392,
        nfcChipId: 'jn8UCSuJK',
        group: {
          id: 3,
          name: 'Mylan Pharmaceuticals Inc.'
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
        id: 36,
        name: 'Margaux MacMenamy',
        saldo: 456,
        nfcChipId: 'kR8FeFup7',
        group: {
          id: 2,
          name: 'A-S Medication Solutions LLC',
          description: 'E.E.S',
          canOverdraw: true
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
        id: 38,
        name: 'Carissa Whitbread',
        saldo: 446,
        nfcChipId: '7Hr4N8G3',
        group: {
          id: 3,
          name: 'Mylan Pharmaceuticals Inc.'
        }
      },
      {
        id: 39,
        name: 'Malia MacElholm',
        saldo: 432,
        nfcChipId: 'pyFmknSCsI2',
        group: {
          id: 9,
          name: 'Pharmacia and Upjohn Company'
        }
      },
      {
        id: 40,
        name: 'Eugenius Odell',
        description: 'BRIGHTER BY NATURE',
        saldo: 477,
        nfcChipId: 'QlRKt1rwd',
        group: {
          id: 2,
          name: 'A-S Medication Solutions LLC',
          description: 'E.E.S',
          canOverdraw: true
        }
      },
      {
        id: 41,
        name: 'Alfy Pietroni',
        description: 'SENSAI CELLULAR PERFORMANCE POWDER FOUNDATION',
        saldo: 428,
        nfcChipId: 'bGK06Cy',
        group: {
          id: 2,
          name: 'A-S Medication Solutions LLC',
          description: 'E.E.S',
          canOverdraw: true
        }
      },
      {
        id: 42,
        name: 'Michele Ondrus',
        saldo: 364,
        nfcChipId: 'QbJuR2kx',
        group: {
          id: 8,
          name: 'Kareway Product, Inc.',
          description: 'Acetaminophen',
          canOverdraw: true
        }
      },
      {
        id: 43,
        name: 'Ulla Risbridge',
        saldo: 447,
        nfcChipId: 'eOPZIlyhkF',
        group: {
          id: 4,
          name: 'Mylan Pharmaceuticals Inc.',
          description: 'Enalapril Maleate and Hydrochlorothiazide'
        }
      },
      {
        id: 44,
        name: 'Hermione Forsaith',
        description: 'Degree',
        saldo: 401,
        nfcChipId: 'D5KMsugN9',
        group: {
          id: 5,
          name: 'REMEDYREPACK INC.',
          description: 'CELEBREX'
        }
      },
      {
        id: 45,
        name: 'Patrice Kigelman',
        description: 'equaline nicotine',
        saldo: 388,
        nfcChipId: 'Q65typehDim',
        group: {
          id: 7,
          name: 'PSS World Medical, Inc.',
          canOverdraw: true
        }
      },
      {
        id: 46,
        name: 'Augy Scraney',
        description: 'Baclofen',
        saldo: 359,
        nfcChipId: 'EQlaU7g',
        group: {
          id: 3,
          name: 'Mylan Pharmaceuticals Inc.'
        }
      },
      {
        id: 47,
        name: 'Massimiliano Fender',
        saldo: 375,
        nfcChipId: 'npzeiAR42qN',
        group: {
          id: 7,
          name: 'PSS World Medical, Inc.',
          canOverdraw: true
        }
      },
      {
        id: 48,
        name: 'Maria Grass',
        description: '4 in 1 Pressed Mineral SPF 15 Light',
        saldo: 468,
        nfcChipId: 'GBnDATCxBL',
        group: {
          id: 5,
          name: 'REMEDYREPACK INC.',
          description: 'CELEBREX'
        }
      },
      {
        id: 49,
        name: 'Francois Dener',
        description: 'topcare cold and flu night time',
        saldo: 415,
        nfcChipId: 'XAsvh8',
        group: {
          id: 4,
          name: 'Mylan Pharmaceuticals Inc.',
          description: 'Enalapril Maleate and Hydrochlorothiazide'
        }
      },
      {
        id: 50,
        name: 'Sebastiano Purselowe',
        saldo: 411,
        nfcChipId: 'udvVRbEIDcR5',
        group: {
          id: 3,
          name: 'Mylan Pharmaceuticals Inc.'
        }
      }
    ]
  }

  searchOnTable() {
    if (this.search) {
      let searchTerm = this.search.toLowerCase()
      this.searched = this.accounts.filter((item) => item.name.toLowerCase().includes(searchTerm))
    }
  }

  deleteAccount(account: Account) {
    this.accounts = this.accounts.reduce(function (prev: Account[], curr: Account) {
      if (curr.id === account.id) {
        return prev
      }
      prev.push(curr)
      return prev
    }, [])

    this.searched = this.accounts
  }

  editRoute(accountId: number) {
    return `/accounts/${accountId}/edit`
  }

  transactionsRoute(accountId: number) {
    return `/accounts/${accountId}/transactions`
  }

  created() {
    this.getAccounts()
    this.searched = this.accounts
  }
}
</script>

<style lang="scss" scoped>
  .md-field {
    max-width: 300px;
  }
</style>

<style lang="scss">
  .md-card.md-app-content {
    height: auto;
  }
</style>
