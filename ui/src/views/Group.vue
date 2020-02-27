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
import Group, { emptyGroup } from '@/data/group'
import Account from '@/data/account'
import axios from 'axios'

@Component({
  components: { AccountList, GroupForm }
})
export default class GroupDetail extends Vue {
  group: Group = emptyGroup
  accounts: Account[] = []

  created() {
    axios.get('/group/' + this.$route.params.id, {
      baseURL: 'http://localhost:8088/v1',
      headers: {
        // @ts-ignore
        Authorization: 'Bearer ' + JSON.parse(localStorage.getItem('jwt')).access_token
      }
    }).then((response) => {
      let resGroup = response.data
      this.group = {
        id: resGroup.id,
        name: resGroup.name,
        description: resGroup.hasOwnProperty('description') ? resGroup.description : '',
        canOverdraw: resGroup.can_overdraw
      }
    }).catch((response) => {
      console.error(response)
    })

    axios.get('/accounts', {
      baseURL: 'http://localhost:8088/v1',
      headers: {
        Authorization: 'Bearer ' + JSON.parse(localStorage.getItem('jwt')).access_token
      },
      params: {
        'group_id': this.$route.params.id
      }
    }).then((response) => {
      this.accounts = response.data.accounts
    }).catch((response) => {
      console.error(response)
    })
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
