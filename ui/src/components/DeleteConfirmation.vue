<template>
  <md-button @click="active = true" class="md-icon-button md-accent">
    <md-icon>delete</md-icon>
    <md-dialog-confirm
      :md-active.sync="active"
      :md-title=title
      :md-content=content
      md-confirm-text="Yes"
      md-cancel-text="No"
      @md-confirm="onConfirm"/>
  </md-button>
</template>

<script lang="ts">
import { Component, Emit, Prop, Vue } from 'vue-property-decorator'

@Component
export default class DeleteConfirmation extends Vue {
  active: boolean = false
  @Prop() title!: string
  @Prop() content!: string

  @Emit('onDelete') onConfirm() {
    return false
  }
}
</script>

<style scoped>
  .delete-btn {
    display: inline-block;
  }
</style>
