<template>
  <v-text-field
    v-bind:value="value"
    @input="$emit('input', $event)"
    :counter="255"
    :rules="[(v) => !!v || 'Name is required']"
    label="Name"
    required
    prepend-icon="person"
  />
</template>

<script lang="ts">
import { Component, Prop, Vue } from 'vue-property-decorator'

@Component
export default class AccountFieldName extends Vue {
  @Prop({
    type: String,
    required: true
  })
  value!: string
}
</script>

<style scoped>

</style>
