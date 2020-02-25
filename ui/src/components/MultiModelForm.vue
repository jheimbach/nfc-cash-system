<template>
  <v-container>
    <v-row>
      <v-col cols="12" md="12">
        <v-text-field
          v-model="len"
          type="number"
          :label="lenLabel"
          required
          min="0"
          :hint="lenHint"
        ></v-text-field>
      </v-col>
    </v-row>
    <v-form v-model="valid[index]" ref="form" class="form" v-for="(model, index) in models" :key="'model_'+index"
            lazy-validation>
      <slot name="form" v-bind:model="model" v-bind:removeModel="removeModel" v-bind:index="index">
        <v-row align="center">
          <v-col cols="12" md="2">
            <account-field-name v-model="model.name"/>
          </v-col>

          <v-col cols="12" md="3">
            <field-description v-model="model.description"/>
          </v-col>

          <v-col cols="12" md="2">
            <account-field-saldo v-model="model.saldo"/>
          </v-col>

          <v-col cols="12" md="2">
            <account-field-group v-model="model.group"/>
          </v-col>

          <v-col cols="12" md="2">
            <account-field-nfc-chip-id v-model="model.nfcChipId"/>
          </v-col>
          <v-col>
            <v-btn icon @click="removeModel(index)">
              <v-icon>delete</v-icon>
            </v-btn>
          </v-col>
        </v-row>
      </slot>
    </v-form>
    <v-row>
      <v-col cols="12" md="12">
        <v-btn :disabled="!isValid()" color="success" class="mr-4" @click="save">{{saveBtnTxt}}</v-btn>
      </v-col>
    </v-row>
    <slot name="snackbar" v-bind:activator="snackbar" v-bind:hideSnackbar="hideSnackbar">
      <v-snackbar :value="snackbar" @input="hideSnackbar">
        Accounts successfully created
        <v-btn color="pink" text @click="hideSnackbar">Close</v-btn>
      </v-snackbar>
    </slot>
  </v-container>
</template>

<script lang="ts">
import { Component, Prop, Vue, Watch, Emit } from 'vue-property-decorator'
import AccountFieldName from '@/components/form/Account/Name.vue'
import AccountFieldNfcChipId from '@/components/form/Account/NfcChipId.vue'
import AccountFieldGroup from '@/components/form/Account/Group.vue'
import AccountFieldSaldo from '@/components/form/Account/Saldo.vue'
import { emptyAccount } from '@/data/account'
import FieldDescription from '@/components/form/Description.vue'

@Component({
  components: { AccountFieldName, FieldDescription, AccountFieldNfcChipId, AccountFieldGroup, AccountFieldSaldo }
})
export default class MultiModelForm extends Vue {
  @Prop({
    type: Array,
    required: true
  })
  models!: any
  valid: boolean[] = []
  snackbar: boolean = false
  nModel: number = 0

  get len() {
    return this.nModel.toString()
  }

  set len(val: any) {
    this.nModel = parseInt(val)
  }

  @Prop({
    type: String,
    default: 'No. of Accounts'
  })
  lenLabel!: string

  @Prop({
    type: Number,
    default: 5
  })
  defaultLen!: number

  @Prop({
    type: String,
    default: 'how many accounts do you want to be created?'
  })
  lenHint!: string

  @Prop({
    type: String,
    default: 'Create Accounts'
  })
  saveBtnTxt!: string

  @Watch('len')
  @Emit('lenUpdated')
  updateAccountLength(): number {
    const nModels = this.nModel - this.models.length
    if (nModels > 0) {
      for (let i = 0; i < nModels; i++) {
        this.valid.push(false)
      }
    }
    if (nModels < 0) {
      this.valid = this.valid.slice(0, nModels)
    }
    return this.nModel
  }

  @Emit('removeModel')
  removeModel(idx: number): number {
    this.len = (this.nModel - 1)
    this.valid.splice(idx, 1)
    return idx
  }

  created() {
    this.len = this.defaultLen
  }

  hideSnackbar() {
    this.snackbar = false
  }

  isValid() {
    for (let i = 0; i < this.valid.length; i++) {
      if (!this.valid[i]) {
        return false
      }
    }
    return true
  }

  save() {
    // @ts-ignore
    const isValid = this.$refs.form.reduce((all, el) => {
      if (!el.validate()) {
        all = false
      }
      return all
    }, true)
    if (isValid) {
      this.snackbar = true
    }
  }
}
</script>

<style scoped lang="scss">
  .form {
    width: 100%;
  }
</style>
