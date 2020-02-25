<template>
  <v-text-field
    v-bind:value="value"
    @input="$emit('input', $event)"
    :counter="20"
    :rules="chipRules"
    label="NFC Chip ID"
    required
    prepend-icon="nfc"
  />
</template>

<script lang="ts">
import { Component, Prop, Vue } from 'vue-property-decorator'

@Component
export default class AccountFieldNfcChipId extends Vue {
  @Prop({
    type: String,
    required: true
  })
  value!: string
  chipRules = [
    (v: string) => !!v || 'NFC Chip ID is required',
    (v: string) => (v || '').length <= 20 || 'NFC Chip ID should be max. 20 characters'
  ]
}
</script>

<style scoped>

</style>
