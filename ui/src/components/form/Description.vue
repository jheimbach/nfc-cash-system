<template>
  <v-textarea
    v-bind:value="value"
    @input="$emit('input', $event)"
    label="Description"
    prepend-icon="subject"
    v-if="textarea"
  />
  <v-text-field
    v-bind:value="value"
    @input="$emit('event', $event)"
    label="Description"
    prepend-icon="subject"
    v-else
  />
</template>

<script lang="ts">
import { Component, Prop, Vue } from 'vue-property-decorator'

@Component
export default class FieldDescription extends Vue {
  @Prop({
    type: String,
    required: true
  })
  value!:string
  @Prop({
    type: Boolean,
    default: false
  })
  textarea!:boolean
}
</script>

<style scoped>

</style>
