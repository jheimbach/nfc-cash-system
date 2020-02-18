<template>
  <div class="account-details">
    <div class="row">
      <h3>ID</h3>
      <p class="detail-content">{{account.id}}</p>
    </div>
    <div class="row">
      <h3>Name</h3>
      <p class="detail-content">{{account.name}}</p>
    </div>
    <div class="row">
      <h3>Description</h3>
      <p class="detail-content">{{account.description}}</p>
    </div>
    <div class="row">
      <h3>Saldo</h3>
      <p class="detail-content">{{account.saldo}}</p>
    </div>
    <div class="row">
      <h3>NFC Chip ID</h3>
      <p class="detail-content">{{account.nfcChipId}}</p>
    </div>
    <div class="row">
      <h3>Group</h3>
      <router-link :to="{name:'group_detail', params:{id: account.group.id}}" class="detail-content">
        {{account.group.name}}
      </router-link>
    </div>
    <!--<div class="row">
      <md-field>
        <label>ID</label>
        <md-input v-model="account.id" readonly></md-input>
      </md-field>
    </div>
    <div class="row">
      <md-field>
        <label>Name</label>
        <md-input v-model="account.name"></md-input>
      </md-field>
    </div>
    <div class="row">
      <md-field>
        <label>Description</label>
        <md-textarea v-model="account.description"></md-textarea>
      </md-field>
    </div>
    <div class="row">
      <md-field>
        <label>Saldo</label>
        <md-input v-model="account.saldo" type="number" readonly></md-input>
      </md-field>
    </div>
    <div class="row">
      <md-field>
        <label>NFC Chip ID</label>
        <md-input v-model="account.nfcChipId"></md-input>
      </md-field>
    </div>
    <div class="row">
      <h3>Group</h3>
      <group-detail-list :group="account.group"/>
    </div>-->
  </div>
</template>

<script lang="ts">
import { Component, Prop, Vue } from 'vue-property-decorator'
import Account from '@/data/account'
import GroupDetailList from '@/components/GroupDetailList.vue'

@Component({
  components: { GroupDetailList }
})
export default class AccountDetailList extends Vue {
  @Prop() account!: Account
}
</script>

<style scoped lang="scss">
  .row {
    margin-bottom: 10px;
    display: flex;
    background: #eee;

    h3 {
      margin: 0;
      width: 150px;
      padding: 5px;
      background: #ddd;
    }

    .detail-content {
      margin: 2px 0;
      padding: 5px;
      font-size: 1.4em;
      flex: 1;
      line-height: 1.4;
    }
  }
</style>
