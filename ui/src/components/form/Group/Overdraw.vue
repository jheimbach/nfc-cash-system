<template>
  <v-switch
    v-bind:value="value"
    @input="$emit('input', $event)"
    label="Can Overdraw" prepend-icon="trending_down"/>
</template>

<script lang="ts">
import { Component, Prop, Vue } from 'vue-property-decorator'

@Component
export default class GroupFieldCanOverdraw extends Vue {
  @Prop({
    type: Boolean,
    required: true,
    default: false
  })
  value!: string
}
</script>

<style scoped>

</style>
