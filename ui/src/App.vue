<template>
  <v-app id="inspire">
    <v-navigation-drawer v-model="drawer" app>
      <v-list dense>
        <v-list-item link v-for="link in navigation" :to="link.to" @click.stop="toggleDrawer" :key="link.title">
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
      <v-app-bar-nav-icon @click.stop="toggleDrawer"/>
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
  <!--<div class="page-container">
     <md-app md-mode="overlap">
       <md-app-toolbar class="md-large md-primary">
         <div class="md-toolbar-row">
           <div class="md-toolbar-section-start">
             <md-button class="md-icon-button" @click="toggleMenu">
               <md-icon>menu</md-icon>
             </md-button>
           </div>
           <div class="md-toolbar-row md-toolbar-offset">
             <span class="md-title">NFC Cash System</span>
           </div>
         </div>
       </md-app-toolbar>
       <md-app-drawer :md-active.sync="menuVisible">
         <md-toolbar class="md-transparent" md-elevation="0">
           <span class="md-title">Navigation</span>
           <div class="md-toolbar-section-end">
             <md-button class="md-icon-button md-dense" @click="toggleMenu">
               <md-icon>keyboard_arrow_left</md-icon>
             </md-button>
           </div>
         </md-toolbar>
         <md-list>
           <md-list-item to="/" exact @click="toggleMenu">
             <md-icon>move_to_inbox</md-icon>
             <span class="md-list-item-text">Home</span>
           </md-list-item>
           <md-list-item to="/accounts" @click="toggleMenu">
             <md-icon>account_box</md-icon>
             <span class="md-list-item-text">Accounts</span>
           </md-list-item>
           <md-list-item to="/groups" @click="toggleMenu">
             <md-icon>group</md-icon>
             <span class="md-list-item-text">Groups</span>
           </md-list-item>
           <md-list-item to="/transactions" @click="toggleMenu">
             <md-icon>account_balance_wallet</md-icon>
             <span class="md-list-item-text">Transactions</span>
           </md-list-item>
         </md-list>
       </md-app-drawer>
       <md-app-content>
         <router-view/>
       </md-app-content>
     </md-app>
   </div>-->
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
      title: 'home',
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

  toggleDrawer() {
    this.drawer = !this.drawer
  }
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
