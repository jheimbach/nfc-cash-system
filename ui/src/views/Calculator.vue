<template>
  <div class="calculator">
    <v-container v-if="calcInput !== '' || entries.length > 0" class="order-container">
      <div ref="entrylist" class="entrylist">
        <v-row v-for="(entry, index) in entries" :key="`entry-${index}`">
          <v-col cols="1" class="calc-entry-method text-right">{{entry.method}}</v-col>
          <v-col class="display-3 text-right" cols="11">
            <span class="calc-entry-text blue-grey--text lighten-5"
                  v-if="entry.text !== entry.value.toString()">{{entry.text}} =
            </span>
            {{entry.value|formatMoney}}
          </v-col>
        </v-row>
      </div>
      <v-row v-if="sum !== 0" class="calc-total-row">
        <v-col cols="1" class=" display-3 text-right">
          Total:
        </v-col>
        <v-col class="display-3 text-right">
          {{sum|formatMoney}}
        </v-col>
      </v-row>
    </v-container>
    <v-container class="order-container empty" v-else>
      <v-row>
        <v-col class="display-4">
          {{emptyText}}
        </v-col>
      </v-row>
    </v-container>
    <v-container>
      <v-row>
        <v-text-field ref="inputField" v-model="calcInput" reverse height="100" class="calc-input-field"
                      @keydown="checkInput" autofocus :error-messages="errMsgs" outlined label="Enter Amount"/>
      </v-row>
    </v-container>
    <v-card>
      <v-card-title>{{account.name}}</v-card-title>
      <v-card-text>
        <v-container>
          <v-row>
            <v-col md="6" cols="12">
              <div class="account-info-header font-weight-bold subtitle-1">
                <v-icon left>subject</v-icon>
                Description
              </div>
              <div class="account-info body-1 description">{{account.description}}</div>
            </v-col>
            <v-col md="4" cols="12">
              <div class="account-info-header font-weight-bold subtitle-1">
                <v-icon left>attach_money</v-icon>
                Saldo
              </div>
              <div class="account-info body-1 saldo">{{account.saldo|formatMoney}}</div>
            </v-col>
            <v-col md="2" cols="12">
              <div class="account-info-header font-weight-bold subtitle-1">
                <v-icon left>group</v-icon>
                Group
              </div>
              <div class="account-info body-1 group">{{account.group.name}}</div>
            </v-col>
          </v-row>
        </v-container>
      </v-card-text>
    </v-card>
    <v-snackbar :value="snackbar" @input="hideSnackbar">
      Transaction for {{account.name}} created
      <v-btn color="pink" text @click="hideSnackbar">Close</v-btn>
    </v-snackbar>
  </div>
</template>

<script lang="ts">
import { Component, Vue, Watch } from 'vue-property-decorator'
import Account from '@/data/account'

enum EntryMethod {
  Plus = '+',
  Minus = '-',
}

interface CalcEntry {
  value: number,
  method: EntryMethod,
  text: string,
}

@Component
export default class Calculator extends Vue {
  account: Account = {
    id: 1,
    name: 'Tim Strupp',
    description: 'Lorem ipsum dolor sit amet, consectetur adipisicing elit.',
    saldo: 123,
    group: {
      id: 1,
      name: 'test Group 1'
    },
    nfcChipId: '123456789'
  }
  entries: CalcEntry[] = []
  sum: number = 0
  calcInput: string = ''
  errMsgs: string[] = []
  commandKeys: string[] = []
  allowedKeys: string[] = []
  snackbar: boolean = false
  emptyText: string = 'Hello! Your order please!'

  created() {
    this.commandKeys = '+-'.split('')
    this.commandKeys.push(...['Enter', 'Backspace', 'Escape'])
    this.allowedKeys = '0123456789.,/*'.split('')
    this.allowedKeys.push(...this.commandKeys)
  }

  @Watch('calcInput')
  validateInput() {
    if (this.calcInput.indexOf('/0') >= 0) {
      this.errMsgs = ['Can not divide by zero']
    } else {
      this.errMsgs = []
    }
  }

  @Watch('entries')
  scrollToBottomOfList() {
    const list = this.$refs.entrylist as HTMLElement
    if (list) {
      this.$nextTick(() => {
        list.scrollTo({ top: list.scrollHeight })
      })
    }
  }

  escapePressed: number = 0

  handleEscape() {
    if (this.escapePressed >= 1) {
      if (this.calcInput === '') {
        this.entries = []
        this.sum = 0
      } else {
        this.calcInput = ''
      }
      this.escapePressed = 0
    } else {
      this.escapePressed++
    }
  }

  checkInput(ev: KeyboardEvent) {
    if (this.allowedKeys.indexOf(ev.key) === -1) {
      ev.preventDefault()
      return
    }
    if (this.commandKeys.indexOf(ev.key) !== -1) {
      if (ev.key === 'Backspace') {
        this.retrieveLastLine(ev)
        return
      }
      if (ev.key === 'Escape') {
        this.handleEscape()
        return
      }
      if (ev.key === 'Enter') {
        this.handleSave()
        return
      }
    }
    if (this.errMsgs.length !== 0) {
      ev.preventDefault()
      return
    }
    switch (ev.key) {
      case '+': {
        this.createCalcEntry(EntryMethod.Plus)
        ev.preventDefault()
        break
      }
      case '-': {
        this.createCalcEntry(EntryMethod.Minus)
        ev.preventDefault()
        break
      }
    }
  }

  retrieveLastLine(ev:KeyboardEvent) {
    if (this.calcInput === '' && this.entries.length !== 0) {
      const lastEntry = this.entries[this.entries.length - 1]
      this.entries = this.entries.slice(0, this.entries.length - 1)
      if (lastEntry.method === EntryMethod.Plus) {
        this.sum = this.sum - lastEntry.value
      } else {
        this.sum = this.sum + lastEntry.value
      }
      this.calcInput = lastEntry.text
      ev.preventDefault()
    }
  }

  createCalcEntry(method: EntryMethod) {
    if (this.calcInput === '') {
      return
    }
    const val = this.parseInputVal(this.calcInput)
    if (method === EntryMethod.Plus) {
      this.sum += val
    } else {
      this.sum -= val
    }
    this.entries.push({
      value: val,
      text: this.calcInput.replace(method, ''),
      method: method
    })
    this.calcInput = ''
  }

  parseInputVal(val: string): number {
    let multiValues = val.split('*')
    if (multiValues.length > 1) {
      return this.multiplyValue(multiValues)
    }
    let divVal = val.split('/')
    if (divVal.length > 1) {
      return this.divideValue(divVal)
    }
    return parseFloat(val)
  }

  multiplyValue(vals: string[]) {
    let mVal = parseFloat(vals[0])
    for (let i = 1; i < vals.length; i++) {
      mVal *= parseFloat(vals[i])
    }
    return mVal
  }

  divideValue(vals: string[]) {
    let dVal = parseFloat(vals[0])
    for (let i = 1; i < vals.length; i++) {
      dVal = dVal / parseFloat(vals[i])
    }
    return dVal
  }

  handleSave() {
    // create entry if calcInput is not empty
    this.createCalcEntry(EntryMethod.Plus)
    setTimeout(() => {
      this.snackbar = true
      this.calcInput = ''
      this.entries = []
      this.sum = 0
    }, 500)
  }

  hideSnackbar() {
    this.snackbar = false
  }
}
</script>

<style scoped lang="scss">
  @import "~vuetify/src/styles/settings/colors";

  .calculator {
    width: 100%;
  }

  .account-info {
    padding-left: 36px;
  }

  .calc-entry-method {
    font-size: 2em;
    align-self: center;
  }

  .calc-entry-text {
    opacity: .2;
    font-size: .7em;
  }

  .calc-input-field {
    font-size: 80px;
  }

  .entrylist {
    height: 50vh;
    overflow-y: auto;
    width: calc(100% + 24px);
    margin-left: -12px;

    .row {
      width: calc(100% - 24px);
    }
  }

  .order-container.empty {
    height: 50vh;
  }

  .calc-total-row {
    margin-top: 20px;
    border-top: 2px solid map-deep-get($blue-grey, 'lighten-4');
  }
</style>
<style lang="scss">
  .v-text-field--outlined .v-label--active {
    transform: translateY(-14px) scale(0.75);
  }
</style>
