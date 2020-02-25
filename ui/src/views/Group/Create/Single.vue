<template>
  <group-form :group="group" @save="save" save-btn-txt="Create group" cancel-btn-txt="Clear form"/>
</template>

<script lang="ts">

import { Component, Vue } from 'vue-property-decorator'
import GroupForm from '@/components/groups/Form.vue'
import Group, { emptyGroup } from '@/data/group'

@Component({
  components: { GroupForm }
})
export default class GroupCreateSingle extends Vue {
  group: Group = Object.assign({}, emptyGroup)

  save() {
    this.group = emptyGroup
  }
}
</script>

<style scoped>

</style>
