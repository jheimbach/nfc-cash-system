<template>
  <v-form ref="form" lazy-validation v-model="valid">
    <v-text-field
      v-model="group.id"
      label="ID"
      disabled
      prepend-icon="vpn_key"
      v-if="group.id !== 0"
    />
    <v-text-field
      v-model="group.name"
      :counter="255"
      :rules="[(v) => !!v || 'Name is required']"
      label="Name"
      required
      prepend-icon="person"
    />
    <v-textarea
      v-model="group.description"
      label="Description"
      prepend-icon="subject"
    />
    <v-switch v-model="group.canOverdraw" label="Can Overdraw" prepend-icon="trending_down"/>

    <v-btn :disabled="!valid" color="success" class="mr-4" @click="save">Save</v-btn>
    <v-btn color="error" class="mr-4" @click="cancel">Clear Changes</v-btn>

    <v-snackbar v-model="snackbar">
      Group {{group.name }} successfully saved
      <v-btn color="pink" text @click="snackbar = false">Close</v-btn>
    </v-snackbar>
  </v-form>
</template>

<script lang="ts">
import { Component, Prop, Vue, Emit } from 'vue-property-decorator'
import Group from '@/data/group'

@Component
export default class GroupForm extends Vue {
  @Prop({
    default: () => {
      return {
        id: 0,
        name: '',
        description: '',
        canOverdraw: false
      }
    },
    required: true,
    type: Object
  })
  group!: Group
  unchangedGroup!: Group
  valid: boolean = false
  snackbar: boolean = false

  @Emit('save')
  save() {
    // @ts-ignore
    if (this.$refs.form.validate()) {
      setTimeout(() => {
        this.snackbar = true
      }, 500)
    }

    return this.group
  }

  @Emit('cancel')
  cancel() {
    this.group.id = this.unchangedGroup.id
    this.group.name = this.unchangedGroup.name
    this.group.description = this.unchangedGroup.description
    this.group.canOverdraw = this.unchangedGroup.canOverdraw
    return this.group
  }

  created() {
    this.unchangedGroup = JSON.parse(JSON.stringify(this.group))
  }
}
</script>

<style scoped>

</style>
