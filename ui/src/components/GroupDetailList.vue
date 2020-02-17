<template>
  <div class="group-detail">
    <div class="row">
      <div class="col field-name">ID:</div>
      <div class="col field-value">{{group.id}}</div>
    </div>
    <div class="row">
      <div class="col field-name">Name:</div>
      <div class="col field-value">{{group.name}}</div>
    </div>
    <div class="row" v-if="showDescription && group.description != null">
      <div class="col field-name">Description:</div>
      <div class="col field-value">{{group.description}}</div>
    </div>
    <div class="row" v-if="showCanOverDraw">
      <div class="col field-name">Can overdraw?:</div>
      <div class="col field-value">{{group.canOverdraw? 'Yes': 'No'}}</div>
    </div>
  </div>
</template>

<script lang="ts">
import { Component, Prop, Vue } from 'vue-property-decorator'
import Group from '../data/group'

@Component
export default class GroupDetailList extends Vue {
  @Prop() showDescription: boolean = false
  @Prop() showCanOverDraw: boolean = false
  @Prop() class: string = ''
  @Prop() group!: Group

  created() {
    if (this.class !== '') {
      this.class = 'group-detail'
    } else {
      this.class = this.class + ' group-detail'
    }
  }
}
</script>

<style scoped lang="scss">
  .group-detail {
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
  }
</style>
