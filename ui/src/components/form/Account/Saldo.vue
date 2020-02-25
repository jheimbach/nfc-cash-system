<template>
  <v-text-field
    v-bind:value="value"
    @input="$emit('input', $event)"
    type="number"
    label="Saldo"
    prepend-icon="attach_money"
    :disabled="disabled"
  />
</template>

<script lang="ts">
import { Component, Prop, Vue } from 'vue-property-decorator'

@Component
export default class AccountFieldSaldo extends Vue {
  @Prop({
    type: Number,
    required: true
  })
  value!: number
  @Prop({
    type: Boolean,
    default: false
  })
  disabled!: boolean
}
</script>

<style scoped>

</style>
