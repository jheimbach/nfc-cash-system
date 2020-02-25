<template>
  <multi-model-form :models="groups"
                    len-label="No. of Groups"
                    len-hint="how many groups do you want to be created?"
                    save-btn-txt="Create Groups"
                    @removeModel="removeGroup"
                    @lenUpdated="updateGroups"
  >
    <template v-slot:form="{model, removeModel, index}">
      <v-row align="center">
        <v-col cols="12" md="3">
          <group-field-name v-model="model.name"/>
        </v-col>
        <v-col cols="12" md="5">
          <field-description v-model="model.description"/>
        </v-col>
        <v-col cols="12" md="3">
          <group-field-can-overdraw v-model="model.canOverdraw"/>
        </v-col>
        <v-col cols="12" md="1">
          <v-btn icon @click="removeModel(index)">
            <v-icon>delete</v-icon>
          </v-btn>
        </v-col>
      </v-row>
    </template>
  </multi-model-form>
</template>

<script lang="ts">

import { Component, Vue } from 'vue-property-decorator'
import MultiModelForm from '@/components/MultiModelForm.vue'
import Group, { emptyGroup } from '@/data/group'
import FieldDescription from '@/components/form/Description.vue'
import GroupFieldCanOverdraw from '@/components/form/Group/Overdraw.vue'
import GroupFieldName from '@/components/form/Group/Name.vue'

@Component({
  components: { MultiModelForm, FieldDescription, GroupFieldName, GroupFieldCanOverdraw }
})
export default class GroupCreateMultiple extends Vue {
  groups: Group[] = []

  updateGroups(len: number) {
    const numberOfAccounts = len - this.groups.length
    if (numberOfAccounts > 0) {
      for (let i = 0; i < numberOfAccounts; i++) {
        this.groups.push(Object.assign({}, emptyGroup))
      }
    }
    if (numberOfAccounts < 0) {
      this.groups = this.groups.slice(0, numberOfAccounts)
    }
  }

  removeGroup(idx: number) {
    this.groups.splice(idx, 1)
  }
}
</script>

<style scoped>

</style>
