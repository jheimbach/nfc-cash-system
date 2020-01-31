<template>
  <div class="page-container">
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
          <md-list-item to="/" exact>
            <md-icon>move_to_inbox</md-icon>
            <span class="md-list-item-text">Home</span>
          </md-list-item>
          <md-list-item to="/accounts">
            <md-icon>account_box</md-icon>
            <span class="md-list-item-text">Accounts</span>
          </md-list-item>
          <md-list-item to="/groups">
            <md-icon>group</md-icon>
            <span class="md-list-item-text">Groups</span>
          </md-list-item>
          <md-list-item to="/transactions">
            <md-icon>account_balance_wallet</md-icon>
            <span class="md-list-item-text">Transactions</span>
          </md-list-item>
        </md-list>
      </md-app-drawer>
      <md-app-content>
        <router-view/>
      </md-app-content>
    </md-app>
  </div>
</template>

<script>
export default {
  name: 'App',
  data: () => ({
    menuVisible: false
  }),
  methods: {
    toggleMenu () {
      this.menuVisible = !this.menuVisible
    }
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
