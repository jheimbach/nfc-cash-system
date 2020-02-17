<template>
  <div class="account-details">
    <div class="title">
      <h2>Account: {{account.name}}</h2>
    </div>
    <div class="row">
      <div class="col field-name">ID:</div>
      <div class="col field-value">{{account.id}}</div>
    </div>
    <div class="row">
      <div class="col field-name">Name:</div>
      <div class="col field-value">{{account.name}}</div>
    </div>
    <div class="row">
      <div class="col field-name">Description:</div>
      <div class="col field-value">{{account.description}}</div>
    </div>
    <div class="row">
      <div class="col field-name">Saldo:</div>
      <div class="col field-value">{{account.saldo}}</div>
    </div>
    <div class="row">
      <div class="col field-name">NFC Chip ID:</div>
      <div class="col field-value">{{account.nfcChipId}}</div>
    </div>
    <div class="group-details">
      <h3>Group</h3>
      <group-detail-list :group="account.group"/>
    </div>
  </div>
</template>

<script lang="ts">
import { Component, Prop, Vue } from 'vue-property-decorator'
import Account from '../data/account'
import GroupDetailList from '@/components/GroupDetailList.vue'

@Component({
  components: { GroupDetailList }
})
export default class AccountDetailList extends Vue {
  @Prop() account!: Account
}
</script>

<style scoped lang="scss">
  .account-details {
    display: flex;
    flex-direction: column;

    .row {
      display: flex;
      flex-direction: row;

      .field-name {
        width: 10%;
        font-weight: bold;
      }

      .field-value {
        flex: 1;
      }
    }

    .title {
      h2 {
        margin-bottom: 10px;
      }
    }

    .group-details {
      h3 {
        margin-bottom: 0;
      }
    }
  }
</style>
