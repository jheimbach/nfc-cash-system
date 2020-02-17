<template>
  <div>
    <md-table v-model="searched" md-sort="name" md-sort-order="asc">
      <md-table-toolbar>
        <div class="md-toolbar-section-start">
          <h1 class="md-title">Transactions</h1>
        </div>

        <md-field md-clearable class="md-toolbar-section-end">
          <md-input placeholder="Search by name..." v-model="search" @input="searchOnTable"/>
        </md-field>
      </md-table-toolbar>

      <md-table-empty-state
        md-label="No transactions found"
        :md-description="`No transactions found for this '${search}' query. Try a different search term.`">
      </md-table-empty-state>

      <md-table-row slot="md-table-row" slot-scope="{ item }">
        <md-table-cell md-label="ID" md-sort-by="id" md-numeric>{{ item.id }}</md-table-cell>
        <md-table-cell md-label="OldSaldo" >{{ item.old_saldo }}</md-table-cell>
        <md-table-cell md-label="NewSaldo" >{{ item.new_saldo }}</md-table-cell>
        <md-table-cell md-label="Amount" >{{ item.amount }}</md-table-cell>
        <md-table-cell md-label="Account" md-sort-by="account.name">{{ item.account.name }}</md-table-cell>
        <md-table-cell md-label="Created" md-sort-by="created">{{ item.created | formatDate }}</md-table-cell>
        <md-table-cell md-label="Actions">
          <md-button :to="accountsRoute(item.account.id)" class="md-icon-button ncs-secondary">
            <md-icon>account_box</md-icon>
          </md-button>
        </md-table-cell>
      </md-table-row>
    </md-table>
  </div>
</template>

<script>
const toLower = text => {
  return text.toString().toLowerCase()
}

const searchByName = (items, term) => {
  if (term) {
    return items.filter(item => toLower(item.name).includes(toLower(term)))
  }

  return items
}

export default {
  name: 'TableSearch',
  data: () => ({
    search: null,
    searched: [],
    groups: [{
      'id': 1,
      'old_saldo': 540,
      'new_saldo': 539,
      'amount': 1,
      'created': '2018-12-10T01:58:06Z',
      'account': {
        'id': 20,
        'name': 'Florida Duesberry',
        'saldo': 495,
        'nfc_chip_id': 'rKlqNQsxt',
        'group': { 'id': 9, 'name': 'Pharmacia and Upjohn Company' }
      }
    }, {
      'id': 2,
      'old_saldo': 540,
      'new_saldo': 533,
      'amount': 7,
      'created': '2018-12-10T04:10:15Z',
      'account': {
        'id': 70,
        'name': 'Teri Trosdall',
        'saldo': 406,
        'nfc_chip_id': 'rihk0T2',
        'group': {
          'id': 10,
          'name': 'Dolgencorp, Inc. (DOLLAR GENERAL & REXALL)',
          'description': 'Allergy Relief',
          'can_overdraw': true
        }
      }
    }, {
      'id': 3,
      'old_saldo': 540,
      'new_saldo': 530,
      'amount': 10,
      'created': '2018-12-10T11:06:01Z',
      'account': {
        'id': 69,
        'name': 'Mortie Wilshire',
        'description': 'Muscle and Joint',
        'saldo': 379,
        'nfc_chip_id': 'WzblXiiz',
        'group': { 'id': 6, 'name': 'H E B', 'description': 'night time' }
      }
    }, {
      'id': 4,
      'old_saldo': 540,
      'new_saldo': 525,
      'amount': 15,
      'created': '2018-12-10T19:19:48Z',
      'account': {
        'id': 88,
        'name': 'Adelle Pigeon',
        'description': 'Acetaminophen',
        'saldo': 424,
        'nfc_chip_id': 'p2smYz0j',
        'group': { 'id': 8, 'name': 'Kareway Product, Inc.', 'description': 'Acetaminophen', 'can_overdraw': true }
      }
    }, {
      'id': 5,
      'old_saldo': 540,
      'new_saldo': 527,
      'amount': 13,
      'created': '2018-12-11T04:11:56Z',
      'account': {
        'id': 89,
        'name': 'Woody Croall',
        'description': 'Voltaren',
        'saldo': 419,
        'nfc_chip_id': 'CldYbsh7',
        'group': {
          'id': 10,
          'name': 'Dolgencorp, Inc. (DOLLAR GENERAL & REXALL)',
          'description': 'Allergy Relief',
          'can_overdraw': true
        }
      }
    }, {
      'id': 6,
      'old_saldo': 540,
      'new_saldo': 537,
      'amount': 3,
      'created': '2018-12-11T06:00:38Z',
      'account': {
        'id': 30,
        'name': 'Rollins Fullard',
        'saldo': 424,
        'nfc_chip_id': '2OKpbD3',
        'group': {
          'id': 10,
          'name': 'Dolgencorp, Inc. (DOLLAR GENERAL & REXALL)',
          'description': 'Allergy Relief',
          'can_overdraw': true
        }
      }
    }, {
      'id': 7,
      'old_saldo': 540,
      'new_saldo': 530,
      'amount': 10,
      'created': '2018-12-11T11:37:26Z',
      'account': {
        'id': 85,
        'name': 'Mehetabel Ratley',
        'saldo': 406,
        'nfc_chip_id': 'LGY2xDU',
        'group': { 'id': 7, 'name': 'PSS World Medical, Inc.', 'can_overdraw': true }
      }
    }, {
      'id': 8,
      'old_saldo': 540,
      'new_saldo': 520,
      'amount': 20,
      'created': '2018-12-11T16:08:04Z',
      'account': {
        'id': 23,
        'name': 'Winona Pebworth',
        'saldo': 440,
        'nfc_chip_id': 'LouoW9Joku',
        'group': { 'id': 9, 'name': 'Pharmacia and Upjohn Company' }
      }
    }, {
      'id': 9,
      'old_saldo': 540,
      'new_saldo': 528,
      'amount': 12,
      'created': '2018-12-11T20:49:28Z',
      'account': {
        'id': 59,
        'name': 'Keith Lafontaine',
        'description': 'Incruse Ellipta',
        'saldo': 480,
        'nfc_chip_id': 'x3813cogAAT',
        'group': { 'id': 3, 'name': 'Mylan Pharmaceuticals Inc.' }
      }
    }, {
      'id': 10,
      'old_saldo': 540,
      'new_saldo': 521,
      'amount': 19,
      'created': '2018-12-12T08:57:22Z',
      'account': {
        'id': 21,
        'name': 'Turner Tutton',
        'saldo': 424,
        'nfc_chip_id': 'LECCoD',
        'group': { 'id': 2, 'name': 'A-S Medication Solutions LLC', 'description': 'E.E.S', 'can_overdraw': true }
      }
    }, {
      'id': 11,
      'old_saldo': 540,
      'new_saldo': 534,
      'amount': 6,
      'created': '2018-12-12T10:15:52Z',
      'account': {
        'id': 96,
        'name': 'Nero Tuffley',
        'description': 'LIPOFEN',
        'saldo': 456,
        'nfc_chip_id': 'qiLEpMFd',
        'group': { 'id': 7, 'name': 'PSS World Medical, Inc.', 'can_overdraw': true }
      }
    }, {
      'id': 12,
      'old_saldo': 530,
      'new_saldo': 516,
      'amount': 14,
      'created': '2018-12-12T21:49:49Z',
      'account': {
        'id': 85,
        'name': 'Mehetabel Ratley',
        'saldo': 406,
        'nfc_chip_id': 'LGY2xDU',
        'group': { 'id': 7, 'name': 'PSS World Medical, Inc.', 'can_overdraw': true }
      }
    }, {
      'id': 13,
      'old_saldo': 540,
      'new_saldo': 535,
      'amount': 5,
      'created': '2018-12-13T03:38:16Z',
      'account': {
        'id': 36,
        'name': 'Margaux MacMenamy',
        'saldo': 456,
        'nfc_chip_id': 'kR8FeFup7',
        'group': { 'id': 2, 'name': 'A-S Medication Solutions LLC', 'description': 'E.E.S', 'can_overdraw': true }
      }
    }, {
      'id': 14,
      'old_saldo': 540,
      'new_saldo': 526,
      'amount': 14,
      'created': '2018-12-13T09:04:55Z',
      'account': {
        'id': 17,
        'name': 'Shanna Ace',
        'description': 'Glipizide',
        'saldo': 466,
        'nfc_chip_id': 'RIL1IVl',
        'group': { 'id': 9, 'name': 'Pharmacia and Upjohn Company' }
      }
    }, {
      'id': 15,
      'old_saldo': 540,
      'new_saldo': 532,
      'amount': 8,
      'created': '2018-12-13T20:20:37Z',
      'account': {
        'id': 99,
        'name': 'Doralyn Sharman',
        'description': 'Rivastigmine Tartrate',
        'saldo': 433,
        'nfc_chip_id': 'ObzBYgybHWR',
        'group': {
          'id': 4,
          'name': 'Mylan Pharmaceuticals Inc.',
          'description': 'Enalapril Maleate and Hydrochlorothiazide'
        }
      }
    }, {
      'id': 16,
      'old_saldo': 540,
      'new_saldo': 538,
      'amount': 2,
      'created': '2018-12-14T01:39:22Z',
      'account': {
        'id': 40,
        'name': 'Eugenius Odell',
        'description': 'BRIGHTER BY NATURE',
        'saldo': 477,
        'nfc_chip_id': 'QlRKt1rwd',
        'group': { 'id': 2, 'name': 'A-S Medication Solutions LLC', 'description': 'E.E.S', 'can_overdraw': true }
      }
    }, {
      'id': 17,
      'old_saldo': 535,
      'new_saldo': 534,
      'amount': 1,
      'created': '2018-12-14T06:36:24Z',
      'account': {
        'id': 36,
        'name': 'Margaux MacMenamy',
        'saldo': 456,
        'nfc_chip_id': 'kR8FeFup7',
        'group': { 'id': 2, 'name': 'A-S Medication Solutions LLC', 'description': 'E.E.S', 'can_overdraw': true }
      }
    }, {
      'id': 18,
      'old_saldo': 540,
      'new_saldo': 536,
      'amount': 4,
      'created': '2018-12-14T11:11:55Z',
      'account': {
        'id': 35,
        'name': 'Geoff Cornejo',
        'saldo': 439,
        'nfc_chip_id': 'Uo57KDmf',
        'group': { 'id': 1, 'name': 'H2O Plus' }
      }
    }, {
      'id': 19,
      'old_saldo': 540,
      'new_saldo': 528,
      'amount': 12,
      'created': '2018-12-14T11:23:46Z',
      'account': {
        'id': 75,
        'name': 'Finlay Woollhead',
        'saldo': 456,
        'nfc_chip_id': 'SBIT7Y2Q',
        'group': {
          'id': 10,
          'name': 'Dolgencorp, Inc. (DOLLAR GENERAL & REXALL)',
          'description': 'Allergy Relief',
          'can_overdraw': true
        }
      }
    }, {
      'id': 20,
      'old_saldo': 540,
      'new_saldo': 539,
      'amount': 1,
      'created': '2018-12-14T13:11:16Z',
      'account': {
        'id': 62,
        'name': 'Cherlyn Kahane',
        'description': 'FOSINOPRIL Na',
        'saldo': 435,
        'nfc_chip_id': '9UK6OkGm3',
        'group': { 'id': 2, 'name': 'A-S Medication Solutions LLC', 'description': 'E.E.S', 'can_overdraw': true }
      }
    }, {
      'id': 21,
      'old_saldo': 540,
      'new_saldo': 534,
      'amount': 6,
      'created': '2018-12-15T02:00:53Z',
      'account': {
        'id': 66,
        'name': 'Fabiano Hablet',
        'saldo': 441,
        'nfc_chip_id': 'AH7OkVugcEhM',
        'group': { 'id': 5, 'name': 'REMEDYREPACK INC.', 'description': 'CELEBREX' }
      }
    }, {
      'id': 22,
      'old_saldo': 540,
      'new_saldo': 525,
      'amount': 15,
      'created': '2018-12-16T16:05:28Z',
      'account': {
        'id': 57,
        'name': 'Edmund Jerams',
        'description': 'Lamotrigine',
        'saldo': 447,
        'nfc_chip_id': 'kCodfKMsEw',
        'group': { 'id': 8, 'name': 'Kareway Product, Inc.', 'description': 'Acetaminophen', 'can_overdraw': true }
      }
    }, {
      'id': 23,
      'old_saldo': 540,
      'new_saldo': 538,
      'amount': 2,
      'created': '2018-12-16T22:11:25Z',
      'account': {
        'id': 86,
        'name': 'Sebastian Sarver',
        'description': 'Folic Acid',
        'saldo': 375,
        'nfc_chip_id': 'Kqt4Z5C',
        'group': { 'id': 8, 'name': 'Kareway Product, Inc.', 'description': 'Acetaminophen', 'can_overdraw': true }
      }
    }, {
      'id': 24,
      'old_saldo': 538,
      'new_saldo': 536,
      'amount': 2,
      'created': '2018-12-16T22:21:11Z',
      'account': {
        'id': 86,
        'name': 'Sebastian Sarver',
        'description': 'Folic Acid',
        'saldo': 375,
        'nfc_chip_id': 'Kqt4Z5C',
        'group': { 'id': 8, 'name': 'Kareway Product, Inc.', 'description': 'Acetaminophen', 'can_overdraw': true }
      }
    }, {
      'id': 25,
      'old_saldo': 540,
      'new_saldo': 531,
      'amount': 9,
      'created': '2018-12-17T02:02:41Z',
      'account': {
        'id': 97,
        'name': 'Rosamond Odo',
        'saldo': 445,
        'nfc_chip_id': 'WtOC94iQ0dM',
        'group': { 'id': 5, 'name': 'REMEDYREPACK INC.', 'description': 'CELEBREX' }
      }
    }, {
      'id': 26,
      'old_saldo': 534,
      'new_saldo': 514,
      'amount': 20,
      'created': '2018-12-17T07:01:18Z',
      'account': {
        'id': 36,
        'name': 'Margaux MacMenamy',
        'saldo': 456,
        'nfc_chip_id': 'kR8FeFup7',
        'group': { 'id': 2, 'name': 'A-S Medication Solutions LLC', 'description': 'E.E.S', 'can_overdraw': true }
      }
    }, {
      'id': 27,
      'old_saldo': 540,
      'new_saldo': 537,
      'amount': 3,
      'created': '2018-12-17T08:02:46Z',
      'account': {
        'id': 47,
        'name': 'Massimiliano Fender',
        'saldo': 375,
        'nfc_chip_id': 'npzeiAR42qN',
        'group': { 'id': 7, 'name': 'PSS World Medical, Inc.', 'can_overdraw': true }
      }
    }, {
      'id': 28,
      'old_saldo': 540,
      'new_saldo': 537,
      'amount': 3,
      'created': '2018-12-17T13:54:27Z',
      'account': {
        'id': 50,
        'name': 'Sebastiano Purselowe',
        'saldo': 411,
        'nfc_chip_id': 'udvVRbEIDcR5',
        'group': { 'id': 3, 'name': 'Mylan Pharmaceuticals Inc.' }
      }
    }, {
      'id': 29,
      'old_saldo': 540,
      'new_saldo': 527,
      'amount': 13,
      'created': '2018-12-17T14:50:49Z',
      'account': {
        'id': 83,
        'name': 'Madlin Readman',
        'description': 'allergy relief',
        'saldo': 355,
        'nfc_chip_id': 'rnVDbdk',
        'group': {
          'id': 4,
          'name': 'Mylan Pharmaceuticals Inc.',
          'description': 'Enalapril Maleate and Hydrochlorothiazide'
        }
      }
    }, {
      'id': 30,
      'old_saldo': 521,
      'new_saldo': 504,
      'amount': 17,
      'created': '2018-12-18T12:41:10Z',
      'account': {
        'id': 21,
        'name': 'Turner Tutton',
        'saldo': 424,
        'nfc_chip_id': 'LECCoD',
        'group': { 'id': 2, 'name': 'A-S Medication Solutions LLC', 'description': 'E.E.S', 'can_overdraw': true }
      }
    }, {
      'id': 31,
      'old_saldo': 540,
      'new_saldo': 523,
      'amount': 17,
      'created': '2018-12-18T16:28:26Z',
      'account': {
        'id': 92,
        'name': 'Hewet Haddy',
        'saldo': 384,
        'nfc_chip_id': 'OqQoW3O',
        'group': { 'id': 5, 'name': 'REMEDYREPACK INC.', 'description': 'CELEBREX' }
      }
    }, {
      'id': 32,
      'old_saldo': 540,
      'new_saldo': 522,
      'amount': 18,
      'created': '2018-12-18T22:12:40Z',
      'account': {
        'id': 65,
        'name': 'Roana Grady',
        'description': 'Thermazene',
        'saldo': 442,
        'nfc_chip_id': 'u1egQQtK',
        'group': { 'id': 1, 'name': 'H2O Plus' }
      }
    }, {
      'id': 33,
      'old_saldo': 520,
      'new_saldo': 504,
      'amount': 16,
      'created': '2018-12-19T04:55:46Z',
      'account': {
        'id': 23,
        'name': 'Winona Pebworth',
        'saldo': 440,
        'nfc_chip_id': 'LouoW9Joku',
        'group': { 'id': 9, 'name': 'Pharmacia and Upjohn Company' }
      }
    }, {
      'id': 34,
      'old_saldo': 522,
      'new_saldo': 512,
      'amount': 10,
      'created': '2018-12-19T13:30:13Z',
      'account': {
        'id': 65,
        'name': 'Roana Grady',
        'description': 'Thermazene',
        'saldo': 442,
        'nfc_chip_id': 'u1egQQtK',
        'group': { 'id': 1, 'name': 'H2O Plus' }
      }
    }, {
      'id': 35,
      'old_saldo': 527,
      'new_saldo': 508,
      'amount': 19,
      'created': '2018-12-19T21:05:07Z',
      'account': {
        'id': 89,
        'name': 'Woody Croall',
        'description': 'Voltaren',
        'saldo': 419,
        'nfc_chip_id': 'CldYbsh7',
        'group': {
          'id': 10,
          'name': 'Dolgencorp, Inc. (DOLLAR GENERAL & REXALL)',
          'description': 'Allergy Relief',
          'can_overdraw': true
        }
      }
    }, {
      'id': 36,
      'old_saldo': 540,
      'new_saldo': 531,
      'amount': 9,
      'created': '2018-12-20T01:57:11Z',
      'account': {
        'id': 4,
        'name': 'Gordy Johnson',
        'saldo': 462,
        'nfc_chip_id': 'KTehLLhT',
        'group': { 'id': 3, 'name': 'Mylan Pharmaceuticals Inc.' }
      }
    }, {
      'id': 37,
      'old_saldo': 537,
      'new_saldo': 526,
      'amount': 11,
      'created': '2018-12-20T10:33:05Z',
      'account': {
        'id': 50,
        'name': 'Sebastiano Purselowe',
        'saldo': 411,
        'nfc_chip_id': 'udvVRbEIDcR5',
        'group': { 'id': 3, 'name': 'Mylan Pharmaceuticals Inc.' }
      }
    }, {
      'id': 38,
      'old_saldo': 540,
      'new_saldo': 529,
      'amount': 11,
      'created': '2018-12-21T01:22:01Z',
      'account': {
        'id': 26,
        'name': 'Stafford Brewin',
        'description': 'Constitutional Enhancer',
        'saldo': 431,
        'nfc_chip_id': 'dLz85N6',
        'group': { 'id': 3, 'name': 'Mylan Pharmaceuticals Inc.' }
      }
    }, {
      'id': 39,
      'old_saldo': 540,
      'new_saldo': 527,
      'amount': 13,
      'created': '2018-12-21T16:50:08Z',
      'account': {
        'id': 45,
        'name': 'Patrice Kigelman',
        'description': 'equaline nicotine',
        'saldo': 388,
        'nfc_chip_id': 'Q65typehDim',
        'group': { 'id': 7, 'name': 'PSS World Medical, Inc.', 'can_overdraw': true }
      }
    }, {
      'id': 40,
      'old_saldo': 540,
      'new_saldo': 520,
      'amount': 20,
      'created': '2018-12-22T01:16:40Z',
      'account': {
        'id': 91,
        'name': 'Fidole Scothorne',
        'description': 'Asprin',
        'saldo': 442,
        'nfc_chip_id': '5RIt0F',
        'group': { 'id': 2, 'name': 'A-S Medication Solutions LLC', 'description': 'E.E.S', 'can_overdraw': true }
      }
    }, {
      'id': 41,
      'old_saldo': 540,
      'new_saldo': 537,
      'amount': 3,
      'created': '2018-12-22T04:35:09Z',
      'account': {
        'id': 27,
        'name': 'Kizzee Pinhorn',
        'description': 'Cysto-Conray II',
        'saldo': 439,
        'nfc_chip_id': 'E3Fnu0oKgBC',
        'group': {
          'id': 4,
          'name': 'Mylan Pharmaceuticals Inc.',
          'description': 'Enalapril Maleate and Hydrochlorothiazide'
        }
      }
    }, {
      'id': 42,
      'old_saldo': 540,
      'new_saldo': 539,
      'amount': 1,
      'created': '2018-12-22T05:30:06Z',
      'account': {
        'id': 13,
        'name': 'Kristin Cicullo',
        'description': 'B.S.C AMPUL',
        'saldo': 373,
        'nfc_chip_id': 'DseWgH8AA1l',
        'group': { 'id': 1, 'name': 'H2O Plus' }
      }
    }, {
      'id': 43,
      'old_saldo': 540,
      'new_saldo': 539,
      'amount': 1,
      'created': '2018-12-22T08:55:53Z',
      'account': {
        'id': 100,
        'name': 'Jedd Wederell',
        'description': 'Omega 3 Targeted Relief',
        'saldo': 484,
        'nfc_chip_id': 'znvwE1VKuoz',
        'group': { 'id': 2, 'name': 'A-S Medication Solutions LLC', 'description': 'E.E.S', 'can_overdraw': true }
      }
    }, {
      'id': 44,
      'old_saldo': 540,
      'new_saldo': 525,
      'amount': 15,
      'created': '2018-12-22T13:29:20Z',
      'account': {
        'id': 55,
        'name': 'Alfonse Jervis',
        'saldo': 379,
        'nfc_chip_id': 'LFQ9NA',
        'group': { 'id': 7, 'name': 'PSS World Medical, Inc.', 'can_overdraw': true }
      }
    }, {
      'id': 45,
      'old_saldo': 540,
      'new_saldo': 520,
      'amount': 20,
      'created': '2018-12-22T14:16:29Z',
      'account': {
        'id': 31,
        'name': 'Wallache Bachelor',
        'description': 'Clindamycin Hydrochloride',
        'saldo': 430,
        'nfc_chip_id': 'htiB3siU',
        'group': { 'id': 9, 'name': 'Pharmacia and Upjohn Company' }
      }
    }, {
      'id': 46,
      'old_saldo': 540,
      'new_saldo': 535,
      'amount': 5,
      'created': '2018-12-23T04:59:56Z',
      'account': {
        'id': 33,
        'name': 'Cassie Praundlin',
        'saldo': 448,
        'nfc_chip_id': 'Y8OSKZcIX9',
        'group': { 'id': 9, 'name': 'Pharmacia and Upjohn Company' }
      }
    }, {
      'id': 47,
      'old_saldo': 504,
      'new_saldo': 503,
      'amount': 1,
      'created': '2018-12-23T15:33:24Z',
      'account': {
        'id': 23,
        'name': 'Winona Pebworth',
        'saldo': 440,
        'nfc_chip_id': 'LouoW9Joku',
        'group': { 'id': 9, 'name': 'Pharmacia and Upjohn Company' }
      }
    }, {
      'id': 48,
      'old_saldo': 527,
      'new_saldo': 509,
      'amount': 18,
      'created': '2018-12-24T04:15:30Z',
      'account': {
        'id': 45,
        'name': 'Patrice Kigelman',
        'description': 'equaline nicotine',
        'saldo': 388,
        'nfc_chip_id': 'Q65typehDim',
        'group': { 'id': 7, 'name': 'PSS World Medical, Inc.', 'can_overdraw': true }
      }
    }, {
      'id': 49,
      'old_saldo': 540,
      'new_saldo': 537,
      'amount': 3,
      'created': '2018-12-24T08:44:05Z',
      'account': {
        'id': 37,
        'name': 'Hugh Cunnell',
        'description': 'Nitrotan',
        'saldo': 466,
        'nfc_chip_id': 'eJZHaCH',
        'group': { 'id': 1, 'name': 'H2O Plus' }
      }
    }, {
      'id': 50,
      'old_saldo': 540,
      'new_saldo': 522,
      'amount': 18,
      'created': '2018-12-24T14:17:37Z',
      'account': {
        'id': 95,
        'name': 'Maurizio Golds',
        'saldo': 393,
        'nfc_chip_id': 'ITrhGELoRES',
        'group': { 'id': 6, 'name': 'H E B', 'description': 'night time' }
      }
    }, {
      'id': 51,
      'old_saldo': 540,
      'new_saldo': 533,
      'amount': 7,
      'created': '2018-12-24T21:08:34Z',
      'account': {
        'id': 52,
        'name': 'Renate Gooding',
        'description': 'Hydroxyzine Hydrochloride',
        'saldo': 342,
        'nfc_chip_id': 'kua3XON',
        'group': { 'id': 2, 'name': 'A-S Medication Solutions LLC', 'description': 'E.E.S', 'can_overdraw': true }
      }
    }, {
      'id': 52,
      'old_saldo': 540,
      'new_saldo': 537,
      'amount': 3,
      'created': '2018-12-25T03:14:52Z',
      'account': {
        'id': 93,
        'name': 'Lena Plummer',
        'saldo': 416,
        'nfc_chip_id': 'Vz8oIu2',
        'group': { 'id': 2, 'name': 'A-S Medication Solutions LLC', 'description': 'E.E.S', 'can_overdraw': true }
      }
    }, {
      'id': 53,
      'old_saldo': 537,
      'new_saldo': 528,
      'amount': 9,
      'created': '2018-12-26T04:42:52Z',
      'account': {
        'id': 30,
        'name': 'Rollins Fullard',
        'saldo': 424,
        'nfc_chip_id': '2OKpbD3',
        'group': {
          'id': 10,
          'name': 'Dolgencorp, Inc. (DOLLAR GENERAL & REXALL)',
          'description': 'Allergy Relief',
          'can_overdraw': true
        }
      }
    }, {
      'id': 54,
      'old_saldo': 540,
      'new_saldo': 536,
      'amount': 4,
      'created': '2018-12-26T12:04:10Z',
      'account': {
        'id': 72,
        'name': 'Jerrold Sincock',
        'description': 'Protex',
        'saldo': 434,
        'nfc_chip_id': 'g5RmkbUEJpzK',
        'group': {
          'id': 10,
          'name': 'Dolgencorp, Inc. (DOLLAR GENERAL & REXALL)',
          'description': 'Allergy Relief',
          'can_overdraw': true
        }
      }
    }, {
      'id': 55,
      'old_saldo': 540,
      'new_saldo': 526,
      'amount': 14,
      'created': '2018-12-26T13:31:15Z',
      'account': {
        'id': 67,
        'name': 'Jeralee Terbeek',
        'description': 'Dr Smiths Diaper Rash',
        'saldo': 443,
        'nfc_chip_id': 'HrLL5F',
        'group': { 'id': 5, 'name': 'REMEDYREPACK INC.', 'description': 'CELEBREX' }
      }
    }, {
      'id': 56,
      'old_saldo': 540,
      'new_saldo': 533,
      'amount': 7,
      'created': '2018-12-26T17:54:52Z',
      'account': {
        'id': 8,
        'name': 'Melisa Josowitz',
        'saldo': 461,
        'nfc_chip_id': 'FsBhsEwr',
        'group': { 'id': 9, 'name': 'Pharmacia and Upjohn Company' }
      }
    }, {
      'id': 57,
      'old_saldo': 540,
      'new_saldo': 534,
      'amount': 6,
      'created': '2018-12-26T22:44:33Z',
      'account': {
        'id': 56,
        'name': 'Perceval Strafen',
        'description': 'ORFADIN',
        'saldo': 442,
        'nfc_chip_id': 'bSyAHkd',
        'group': { 'id': 5, 'name': 'REMEDYREPACK INC.', 'description': 'CELEBREX' }
      }
    }, {
      'id': 58,
      'old_saldo': 540,
      'new_saldo': 524,
      'amount': 16,
      'created': '2018-12-27T09:57:01Z',
      'account': {
        'id': 90,
        'name': 'Dorey Netherclift',
        'saldo': 489,
        'nfc_chip_id': 'wMCzJKBbXB',
        'group': { 'id': 6, 'name': 'H E B', 'description': 'night time' }
      }
    }, {
      'id': 59,
      'old_saldo': 535,
      'new_saldo': 533,
      'amount': 2,
      'created': '2018-12-27T13:59:40Z',
      'account': {
        'id': 33,
        'name': 'Cassie Praundlin',
        'saldo': 448,
        'nfc_chip_id': 'Y8OSKZcIX9',
        'group': { 'id': 9, 'name': 'Pharmacia and Upjohn Company' }
      }
    }, {
      'id': 60,
      'old_saldo': 512,
      'new_saldo': 511,
      'amount': 1,
      'created': '2018-12-28T01:24:06Z',
      'account': {
        'id': 65,
        'name': 'Roana Grady',
        'description': 'Thermazene',
        'saldo': 442,
        'nfc_chip_id': 'u1egQQtK',
        'group': { 'id': 1, 'name': 'H2O Plus' }
      }
    }, {
      'id': 61,
      'old_saldo': 540,
      'new_saldo': 531,
      'amount': 9,
      'created': '2018-12-28T13:59:28Z',
      'account': {
        'id': 79,
        'name': 'Kaela Crosser',
        'saldo': 458,
        'nfc_chip_id': 'PtjPRh',
        'group': { 'id': 1, 'name': 'H2O Plus' }
      }
    }, {
      'id': 62,
      'old_saldo': 540,
      'new_saldo': 527,
      'amount': 13,
      'created': '2018-12-29T16:28:42Z',
      'account': {
        'id': 1,
        'name': 'Laverne Blackstock',
        'description': 'Itchy Eye',
        'saldo': 436,
        'nfc_chip_id': 'Hv8mnajqzIKO',
        'group': { 'id': 7, 'name': 'PSS World Medical, Inc.', 'can_overdraw': true }
      }
    }, {
      'id': 63,
      'old_saldo': 540,
      'new_saldo': 520,
      'amount': 20,
      'created': '2018-12-29T19:06:24Z',
      'account': {
        'id': 60,
        'name': 'Casey Clapton',
        'saldo': 414,
        'nfc_chip_id': 'Nt6MCa62z',
        'group': { 'id': 7, 'name': 'PSS World Medical, Inc.', 'can_overdraw': true }
      }
    }, {
      'id': 64,
      'old_saldo': 540,
      'new_saldo': 526,
      'amount': 14,
      'created': '2018-12-29T19:17:43Z',
      'account': {
        'id': 54,
        'name': 'Conrad Rodenburgh',
        'description': 'Methocarbamol',
        'saldo': 449,
        'nfc_chip_id': 'z6EZ749',
        'group': { 'id': 5, 'name': 'REMEDYREPACK INC.', 'description': 'CELEBREX' }
      }
    }, {
      'id': 65,
      'old_saldo': 540,
      'new_saldo': 523,
      'amount': 17,
      'created': '2018-12-30T03:15:57Z',
      'account': {
        'id': 58,
        'name': 'Quinton Howieson',
        'description': 'Ciprofloxacin',
        'saldo': 397,
        'nfc_chip_id': 'FI791LwJnq',
        'group': { 'id': 7, 'name': 'PSS World Medical, Inc.', 'can_overdraw': true }
      }
    }, {
      'id': 66,
      'old_saldo': 540,
      'new_saldo': 529,
      'amount': 11,
      'created': '2018-12-30T13:05:22Z',
      'account': {
        'id': 42,
        'name': 'Michele Ondrus',
        'saldo': 364,
        'nfc_chip_id': 'QbJuR2kx',
        'group': { 'id': 8, 'name': 'Kareway Product, Inc.', 'description': 'Acetaminophen', 'can_overdraw': true }
      }
    }, {
      'id': 67,
      'old_saldo': 540,
      'new_saldo': 521,
      'amount': 19,
      'created': '2018-12-30T15:03:54Z',
      'account': {
        'id': 34,
        'name': 'Alida Burkert',
        'saldo': 392,
        'nfc_chip_id': 'jn8UCSuJK',
        'group': { 'id': 3, 'name': 'Mylan Pharmaceuticals Inc.' }
      }
    }, {
      'id': 68,
      'old_saldo': 540,
      'new_saldo': 530,
      'amount': 10,
      'created': '2018-12-30T15:06:53Z',
      'account': {
        'id': 49,
        'name': 'Francois Dener',
        'description': 'topcare cold and flu night time',
        'saldo': 415,
        'nfc_chip_id': 'XAsvh8',
        'group': {
          'id': 4,
          'name': 'Mylan Pharmaceuticals Inc.',
          'description': 'Enalapril Maleate and Hydrochlorothiazide'
        }
      }
    }, {
      'id': 69,
      'old_saldo': 533,
      'new_saldo': 515,
      'amount': 18,
      'created': '2018-12-30T16:41:00Z',
      'account': {
        'id': 52,
        'name': 'Renate Gooding',
        'description': 'Hydroxyzine Hydrochloride',
        'saldo': 342,
        'nfc_chip_id': 'kua3XON',
        'group': { 'id': 2, 'name': 'A-S Medication Solutions LLC', 'description': 'E.E.S', 'can_overdraw': true }
      }
    }, {
      'id': 70,
      'old_saldo': 540,
      'new_saldo': 535,
      'amount': 5,
      'created': '2018-12-31T11:35:20Z',
      'account': {
        'id': 11,
        'name': 'Lynea Habberjam',
        'saldo': 449,
        'nfc_chip_id': 'uSe3Sj',
        'group': { 'id': 9, 'name': 'Pharmacia and Upjohn Company' }
      }
    }, {
      'id': 71,
      'old_saldo': 526,
      'new_saldo': 518,
      'amount': 8,
      'created': '2019-01-01T03:58:21Z',
      'account': {
        'id': 54,
        'name': 'Conrad Rodenburgh',
        'description': 'Methocarbamol',
        'saldo': 449,
        'nfc_chip_id': 'z6EZ749',
        'group': { 'id': 5, 'name': 'REMEDYREPACK INC.', 'description': 'CELEBREX' }
      }
    }, {
      'id': 72,
      'old_saldo': 540,
      'new_saldo': 535,
      'amount': 5,
      'created': '2019-01-01T04:34:16Z',
      'account': {
        'id': 15,
        'name': 'Virgina Stairmond',
        'saldo': 400,
        'nfc_chip_id': 'nt9th6L7eD',
        'group': { 'id': 9, 'name': 'Pharmacia and Upjohn Company' }
      }
    }, {
      'id': 73,
      'old_saldo': 540,
      'new_saldo': 534,
      'amount': 6,
      'created': '2019-01-01T11:44:09Z',
      'account': {
        'id': 9,
        'name': 'Matias Beininck',
        'saldo': 466,
        'nfc_chip_id': '0E9wTFbJ',
        'group': { 'id': 8, 'name': 'Kareway Product, Inc.', 'description': 'Acetaminophen', 'can_overdraw': true }
      }
    }, {
      'id': 74,
      'old_saldo': 531,
      'new_saldo': 514,
      'amount': 17,
      'created': '2019-01-02T04:27:43Z',
      'account': {
        'id': 97,
        'name': 'Rosamond Odo',
        'saldo': 445,
        'nfc_chip_id': 'WtOC94iQ0dM',
        'group': { 'id': 5, 'name': 'REMEDYREPACK INC.', 'description': 'CELEBREX' }
      }
    }, {
      'id': 75,
      'old_saldo': 540,
      'new_saldo': 538,
      'amount': 2,
      'created': '2019-01-02T15:09:12Z',
      'account': {
        'id': 74,
        'name': 'Mendel Frantsev',
        'saldo': 412,
        'nfc_chip_id': 'zDpiOdoYU',
        'group': { 'id': 7, 'name': 'PSS World Medical, Inc.', 'can_overdraw': true }
      }
    }, {
      'id': 76,
      'old_saldo': 537,
      'new_saldo': 528,
      'amount': 9,
      'created': '2019-01-02T15:51:04Z',
      'account': {
        'id': 47,
        'name': 'Massimiliano Fender',
        'saldo': 375,
        'nfc_chip_id': 'npzeiAR42qN',
        'group': { 'id': 7, 'name': 'PSS World Medical, Inc.', 'can_overdraw': true }
      }
    }, {
      'id': 77,
      'old_saldo': 540,
      'new_saldo': 525,
      'amount': 15,
      'created': '2019-01-02T20:36:52Z',
      'account': {
        'id': 48,
        'name': 'Maria Grass',
        'description': '4 in 1 Pressed Mineral SPF 15 Light',
        'saldo': 468,
        'nfc_chip_id': 'GBnDATCxBL',
        'group': { 'id': 5, 'name': 'REMEDYREPACK INC.', 'description': 'CELEBREX' }
      }
    }, {
      'id': 78,
      'old_saldo': 533,
      'new_saldo': 518,
      'amount': 15,
      'created': '2019-01-02T21:38:21Z',
      'account': {
        'id': 8,
        'name': 'Melisa Josowitz',
        'saldo': 461,
        'nfc_chip_id': 'FsBhsEwr',
        'group': { 'id': 9, 'name': 'Pharmacia and Upjohn Company' }
      }
    }, {
      'id': 79,
      'old_saldo': 540,
      'new_saldo': 536,
      'amount': 4,
      'created': '2019-01-03T03:37:04Z',
      'account': {
        'id': 3,
        'name': 'Winnie Rennolds',
        'description': 'Ofloxacin',
        'saldo': 436,
        'nfc_chip_id': 'ofzGN0eS2K',
        'group': { 'id': 6, 'name': 'H E B', 'description': 'night time' }
      }
    }, {
      'id': 80,
      'old_saldo': 540,
      'new_saldo': 537,
      'amount': 3,
      'created': '2019-01-03T07:11:55Z',
      'account': {
        'id': 51,
        'name': 'Helenelizabeth Aleksidze',
        'description': 'CULTIVATED OATS POLLEN',
        'saldo': 504,
        'nfc_chip_id': 'XER4ZcTHbj',
        'group': { 'id': 6, 'name': 'H E B', 'description': 'night time' }
      }
    }, {
      'id': 81,
      'old_saldo': 536,
      'new_saldo': 523,
      'amount': 13,
      'created': '2019-01-03T23:02:23Z',
      'account': {
        'id': 3,
        'name': 'Winnie Rennolds',
        'description': 'Ofloxacin',
        'saldo': 436,
        'nfc_chip_id': 'ofzGN0eS2K',
        'group': { 'id': 6, 'name': 'H E B', 'description': 'night time' }
      }
    }, {
      'id': 82,
      'old_saldo': 540,
      'new_saldo': 520,
      'amount': 20,
      'created': '2019-01-04T05:48:09Z',
      'account': {
        'id': 87,
        'name': 'Jordan Yellep',
        'saldo': 491,
        'nfc_chip_id': 'ZWuQP61Fb',
        'group': { 'id': 2, 'name': 'A-S Medication Solutions LLC', 'description': 'E.E.S', 'can_overdraw': true }
      }
    }, {
      'id': 83,
      'old_saldo': 534,
      'new_saldo': 518,
      'amount': 16,
      'created': '2019-01-04T14:43:50Z',
      'account': {
        'id': 66,
        'name': 'Fabiano Hablet',
        'saldo': 441,
        'nfc_chip_id': 'AH7OkVugcEhM',
        'group': { 'id': 5, 'name': 'REMEDYREPACK INC.', 'description': 'CELEBREX' }
      }
    }, {
      'id': 84,
      'old_saldo': 540,
      'new_saldo': 529,
      'amount': 11,
      'created': '2019-01-04T22:15:21Z',
      'account': {
        'id': 81,
        'name': 'Sylvester Faull',
        'saldo': 414,
        'nfc_chip_id': 'u6DLWJWd68',
        'group': { 'id': 5, 'name': 'REMEDYREPACK INC.', 'description': 'CELEBREX' }
      }
    }, {
      'id': 85,
      'old_saldo': 539,
      'new_saldo': 521,
      'amount': 18,
      'created': '2019-01-05T18:25:07Z',
      'account': {
        'id': 62,
        'name': 'Cherlyn Kahane',
        'description': 'FOSINOPRIL Na',
        'saldo': 435,
        'nfc_chip_id': '9UK6OkGm3',
        'group': { 'id': 2, 'name': 'A-S Medication Solutions LLC', 'description': 'E.E.S', 'can_overdraw': true }
      }
    }, {
      'id': 86,
      'old_saldo': 540,
      'new_saldo': 524,
      'amount': 16,
      'created': '2019-01-05T20:47:53Z',
      'account': {
        'id': 24,
        'name': 'Horace Barnewell',
        'description': 'Trout',
        'saldo': 407,
        'nfc_chip_id': 'tbaoWaYXbc',
        'group': { 'id': 6, 'name': 'H E B', 'description': 'night time' }
      }
    }, {
      'id': 87,
      'old_saldo': 540,
      'new_saldo': 530,
      'amount': 10,
      'created': '2019-01-06T09:50:38Z',
      'account': {
        'id': 41,
        'name': 'Alfy Pietroni',
        'description': 'SENSAI CELLULAR PERFORMANCE POWDER FOUNDATION',
        'saldo': 428,
        'nfc_chip_id': 'bGK06Cy',
        'group': { 'id': 2, 'name': 'A-S Medication Solutions LLC', 'description': 'E.E.S', 'can_overdraw': true }
      }
    }, {
      'id': 88,
      'old_saldo': 540,
      'new_saldo': 521,
      'amount': 19,
      'created': '2019-01-06T15:36:39Z',
      'account': {
        'id': 80,
        'name': 'Bud Sinderland',
        'description': 'Russian Olive',
        'saldo': 492,
        'nfc_chip_id': '2TuFZZhDnm7X',
        'group': { 'id': 3, 'name': 'Mylan Pharmaceuticals Inc.' }
      }
    }, {
      'id': 89,
      'old_saldo': 540,
      'new_saldo': 525,
      'amount': 15,
      'created': '2019-01-07T07:40:30Z',
      'account': {
        'id': 71,
        'name': 'Ronna Farron',
        'saldo': 424,
        'nfc_chip_id': 'Rag5Xj9',
        'group': {
          'id': 10,
          'name': 'Dolgencorp, Inc. (DOLLAR GENERAL & REXALL)',
          'description': 'Allergy Relief',
          'can_overdraw': true
        }
      }
    }, {
      'id': 90,
      'old_saldo': 528,
      'new_saldo': 517,
      'amount': 11,
      'created': '2019-01-07T16:38:47Z',
      'account': {
        'id': 47,
        'name': 'Massimiliano Fender',
        'saldo': 375,
        'nfc_chip_id': 'npzeiAR42qN',
        'group': { 'id': 7, 'name': 'PSS World Medical, Inc.', 'can_overdraw': true }
      }
    }, {
      'id': 91,
      'old_saldo': 540,
      'new_saldo': 525,
      'amount': 15,
      'created': '2019-01-07T20:55:51Z',
      'account': {
        'id': 16,
        'name': 'Rubi Howey',
        'description': 'Perrigo Hydroquinone',
        'saldo': 508,
        'nfc_chip_id': 'aagVbL',
        'group': {
          'id': 4,
          'name': 'Mylan Pharmaceuticals Inc.',
          'description': 'Enalapril Maleate and Hydrochlorothiazide'
        }
      }
    }, {
      'id': 92,
      'old_saldo': 531,
      'new_saldo': 524,
      'amount': 7,
      'created': '2019-01-07T22:39:53Z',
      'account': {
        'id': 4,
        'name': 'Gordy Johnson',
        'saldo': 462,
        'nfc_chip_id': 'KTehLLhT',
        'group': { 'id': 3, 'name': 'Mylan Pharmaceuticals Inc.' }
      }
    }, {
      'id': 93,
      'old_saldo': 529,
      'new_saldo': 510,
      'amount': 19,
      'created': '2019-01-07T23:20:38Z',
      'account': {
        'id': 81,
        'name': 'Sylvester Faull',
        'saldo': 414,
        'nfc_chip_id': 'u6DLWJWd68',
        'group': { 'id': 5, 'name': 'REMEDYREPACK INC.', 'description': 'CELEBREX' }
      }
    }, {
      'id': 94,
      'old_saldo': 540,
      'new_saldo': 530,
      'amount': 10,
      'created': '2019-01-08T03:46:39Z',
      'account': {
        'id': 84,
        'name': 'Desmond Scheffel',
        'description': 'Wingscale',
        'saldo': 416,
        'nfc_chip_id': 'pki3p7Aia2yc',
        'group': {
          'id': 4,
          'name': 'Mylan Pharmaceuticals Inc.',
          'description': 'Enalapril Maleate and Hydrochlorothiazide'
        }
      }
    }, {
      'id': 95,
      'old_saldo': 530,
      'new_saldo': 527,
      'amount': 3,
      'created': '2019-01-08T20:01:46Z',
      'account': {
        'id': 69,
        'name': 'Mortie Wilshire',
        'description': 'Muscle and Joint',
        'saldo': 379,
        'nfc_chip_id': 'WzblXiiz',
        'group': { 'id': 6, 'name': 'H E B', 'description': 'night time' }
      }
    }, {
      'id': 96,
      'old_saldo': 540,
      'new_saldo': 535,
      'amount': 5,
      'created': '2019-01-08T21:24:13Z',
      'account': {
        'id': 63,
        'name': 'Dinah Bolan',
        'saldo': 423,
        'nfc_chip_id': 'VGEmEmFNzhV',
        'group': { 'id': 5, 'name': 'REMEDYREPACK INC.', 'description': 'CELEBREX' }
      }
    }, {
      'id': 97,
      'old_saldo': 537,
      'new_saldo': 535,
      'amount': 2,
      'created': '2019-01-09T00:35:15Z',
      'account': {
        'id': 37,
        'name': 'Hugh Cunnell',
        'description': 'Nitrotan',
        'saldo': 466,
        'nfc_chip_id': 'eJZHaCH',
        'group': { 'id': 1, 'name': 'H2O Plus' }
      }
    }, {
      'id': 98,
      'old_saldo': 540,
      'new_saldo': 528,
      'amount': 12,
      'created': '2019-01-10T23:28:52Z',
      'account': {
        'id': 98,
        'name': 'Isacco Serrier',
        'saldo': 412,
        'nfc_chip_id': 'ZpRsq1DZ',
        'group': { 'id': 3, 'name': 'Mylan Pharmaceuticals Inc.' }
      }
    }, {
      'id': 99,
      'old_saldo': 533,
      'new_saldo': 529,
      'amount': 4,
      'created': '2019-01-11T20:43:05Z',
      'account': {
        'id': 70,
        'name': 'Teri Trosdall',
        'saldo': 406,
        'nfc_chip_id': 'rihk0T2',
        'group': {
          'id': 10,
          'name': 'Dolgencorp, Inc. (DOLLAR GENERAL & REXALL)',
          'description': 'Allergy Relief',
          'can_overdraw': true
        }
      }
    }, {
      'id': 100,
      'old_saldo': 526,
      'new_saldo': 523,
      'amount': 3,
      'created': '2019-01-12T17:28:31Z',
      'account': {
        'id': 17,
        'name': 'Shanna Ace',
        'description': 'Glipizide',
        'saldo': 466,
        'nfc_chip_id': 'RIL1IVl',
        'group': { 'id': 9, 'name': 'Pharmacia and Upjohn Company' }
      }
    }
    ]
  }),
  methods: {
    searchOnTable () {
      this.searched = searchByName(this.groups, this.search)
    },
    accountsRoute (groupId) {
      return `/accounts/details/${groupId}`
    }
  },
  created () {
    this.searched = this.groups
  }
}
</script>

<style lang="scss" scoped>
  .md-field {
    max-width: 300px;
  }
</style>
