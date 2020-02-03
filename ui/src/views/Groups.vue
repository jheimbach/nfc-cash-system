<template>
  <div>
    <md-table v-model="searched" md-sort="name" md-sort-order="asc">
      <md-table-toolbar>
        <div class="md-toolbar-section-start">
          <h1 class="md-title">Groups</h1>
        </div>

        <md-field md-clearable class="md-toolbar-section-end">
          <md-input placeholder="Search by name..." v-model="search" @input="searchOnTable"/>
        </md-field>
      </md-table-toolbar>

      <md-table-empty-state
        md-label="No groups found"
        :md-description="`No groups found for this '${search}' query. Try a different search term.`">
      </md-table-empty-state>

      <md-table-row slot="md-table-row" slot-scope="{ item }">
        <md-table-cell md-label="ID" md-sort-by="id" md-numeric>{{ item.id }}</md-table-cell>
        <md-table-cell md-label="Name" md-sort-by="name">{{ item.name }}</md-table-cell>
        <md-table-cell md-label="Descripion" md-sort-by="description">{{ item.description }}</md-table-cell>
        <md-table-cell md-label="Can overdraw" md-sort-by="can_overdraw">{{ item.can_overdraw ? 'Yes' :'No' }}</md-table-cell>
        <md-table-cell md-label="Actions">
          <md-button :to="accountsRoute(item.id)" class="md-icon-button ncs-secondary">
            <md-icon>account_box</md-icon>
          </md-button>
          <md-button :to="editRoute(item.id)" class="md-icon-button md-primary">
            <md-icon>edit</md-icon>
          </md-button>
          <md-button :to="deleteRoute(item.id)" class="md-icon-button md-accent">
            <md-icon>delete</md-icon>
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
    groups: [
      {
        'id': 1,
        'name': 'H2O Plus'
      },
      {
        'id': 2,
        'name': 'A-S Medication Solutions LLC',
        'description': 'E.E.S',
        'can_overdraw': true
      },
      {
        'id': 3,
        'name': 'Mylan Pharmaceuticals Inc.'
      },
      {
        'id': 4,
        'name': 'Mylan Pharmaceuticals Inc.',
        'description': 'Enalapril Maleate and Hydrochlorothiazide'
      },
      {
        'id': 5,
        'name': 'REMEDYREPACK INC.',
        'description': 'CELEBREX'
      },
      {
        'id': 6,
        'name': 'H E B',
        'description': 'night time'
      },
      {
        'id': 7,
        'name': 'PSS World Medical, Inc.',
        'can_overdraw': true
      },
      {
        'id': 8,
        'name': 'Kareway Product, Inc.',
        'description': 'Acetaminophen',
        'can_overdraw': true
      },
      {
        'id': 9,
        'name': 'Pharmacia and Upjohn Company'
      },
      {
        'id': 10,
        'name': 'Dolgencorp, Inc. (DOLLAR GENERAL & REXALL)',
        'description': 'Allergy Relief',
        'can_overdraw': true
      }
    ]
  }),
  methods: {
    searchOnTable () {
      this.searched = searchByName(this.groups, this.search)
    },
    editRoute (groupId) {
      return `/groups/${groupId}/edit`
    },
    accountsRoute (groupId) {
      return `/groups/${groupId}/accounts`
    },
    deleteRoute (groupId) {
      return `/groups/${groupId}/delete`
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
