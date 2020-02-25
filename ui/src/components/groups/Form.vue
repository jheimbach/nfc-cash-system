<template>
  <v-form ref="form" lazy-validation v-model="valid">
    <v-text-field
      v-model="group.id"
      label="ID"
      disabled
      prepend-icon="vpn_key"
      v-if="group.id !== 0"
    />
    <group-field-name v-model="group.name"/>
    <field-description v-model="group.description" textarea/>
    <group-field-can-overdraw v-model="group.canOverdraw"/>

    <v-btn :disabled="!valid" color="success" class="mr-4" @click="save">{{saveBtnTxt}}</v-btn>
    <v-btn color="error" class="mr-4" @click="cancel">{{cancelBtnTxt}}</v-btn>

    <v-snackbar v-model="snackbar">
      Group {{group.name }} successfully saved
      <v-btn color="pink" text @click="snackbar = false">Close</v-btn>
    </v-snackbar>
  </v-form>
</template>

<script lang="ts">
import { Component, Prop, Vue, Emit } from 'vue-property-decorator'
import Group from '@/data/group'
import GroupFieldName from '@/components/form/Group/Name.vue'
import GroupFieldCanOverdraw from '@/components/form/Group/Overdraw.vue'
import FieldDescription from '@/components/form/Description.vue'

@Component({
  components: { GroupFieldName, GroupFieldCanOverdraw, FieldDescription }
})
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

  @Prop({
    default: 'save',
    type: String
  })
  saveBtnTxt!: string

  @Prop({
    default: 'Cancel Changes',
    type: String
  })
  cancelBtnTxt!: string

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
