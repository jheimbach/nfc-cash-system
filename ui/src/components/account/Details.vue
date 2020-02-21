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
      <router-link :to="{name:'group', params:{id: account.group.id}}" class="detail-content">
        {{account.group.name}}
      </router-link>
    </div>
    <!---->
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
    display: flex;
    line-height: 1.4;

    &:last-of-type {
      h3 {
        border-bottom: none;
      }
    }

    h3 {
      margin: 0;
      width: 150px;
      padding: 5px;
      font-size: 1em;
      background: rgba(0, 0, 0, 0.12);
      border-bottom: 1px solid rgba(0, 0, 0, 0.55);
    }

    .detail-content {
      flex: 1;
      padding: 5px;
      margin: 2px 0;
    }
  }
</style>
