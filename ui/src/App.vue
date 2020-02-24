<template>
  <v-app id="inspire">
    <v-navigation-drawer v-model="drawer" app>
      <v-list dense>
        <v-list-item v-for="(link, index) in navigation" :to="link.to" :key="link.title" :tabindex="index" exact>
          <v-list-item-action>
            <v-icon>{{link.icon}}</v-icon>
          </v-list-item-action>
          <v-list-item-content>
            <v-list-item-title>{{link.title}}</v-list-item-title>
          </v-list-item-content>
        </v-list-item>
      </v-list>
    </v-navigation-drawer>
    <v-app-bar app color="primary" dark>
      <v-app-bar-nav-icon @click.stop="drawer = !drawer"/>
      <v-toolbar-title>NFC Cash System</v-toolbar-title>
    </v-app-bar>
    <v-content>
      <v-container>
        <v-row>
          <router-view/>
        </v-row>
      </v-container>
    </v-content>
    <v-footer color="secondary" app>
      <span class="white--text">&copy; 2020</span>
    </v-footer>
  </v-app>
</template>

<script lang="ts">
import Vue from 'vue'
import { Component, Prop } from 'vue-property-decorator'
import { Location } from 'vue-router'

interface DrawerLink {
  icon: string,
  title: string,
  to: Location,
}

@Component
export default class App extends Vue {
  @Prop({
    type: String
  })
  source!: string
  drawer: boolean = false
  navigation: DrawerLink[] = [
    {
      icon: 'move_to_inbox',
      title: 'Home',
      to: {
        name: 'home'
      }
    },
    {
      icon: 'account_box',
      title: 'Accounts',
      to: { name: 'accounts' }
    },
    {
      icon: 'group',
      title: 'Groups',
      to: { name: 'groups' }
    },
    {
      icon: 'account_balance_wallet',
      title: 'Transactions',
      to: { name: 'transactions' }
    }]
}

</script>

<style lang="scss">
  .md-app {
    min-height: 100vh;
    border: 1px solid rgba(#000, .12);
  }

  .md-app-toolbar {
    height: 196px;
  }

  .md-app-container {
    overflow: hidden;
  }
</style>
