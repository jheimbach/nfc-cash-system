<template>
  <account-form :account="{
        id: 0,
        name: '',
        description: '',
        saldo: 0,
        nfcChipId: '',
        group: null
      }" @saved="saved"/>
</template>

<script lang="ts">
import { Component, Vue } from 'vue-property-decorator'
import AccountForm from '../components/account/Form.vue'

@Component({
  components: { AccountForm }
})

export default class TestView extends Vue {
  saved(account: Account) {
    console.log(account)
    console.log('saved')
  }
}
</script>

<style scoped>

</style>
