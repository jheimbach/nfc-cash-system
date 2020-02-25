import Vue from 'vue'
import VueRouter, { RouteConfig, RouterOptions } from 'vue-router'
import Home from '../views/Home.vue'
import Accounts from '@/views/Accounts.vue'
import AccountDetail from '@/views/Account.vue'
import GroupDetail from '@/views/Group.vue'

Vue.use(VueRouter)

const routes: RouteConfig[] = [
  {
    path: '/',
    name: 'home',
    component: Home
  },
  {
    path: '/accounts',
    name: 'accounts',
    component: Accounts
  },
  {
    path: '/account/:id',
    name: 'account',
    component: AccountDetail
  },
  {
    path: '/groups',
    name: 'groups',
    component: () => import(/* webpackChunkName: "groups" */ '../views/Groups.vue')
  },
  {
    path: '/group/:id',
    name: 'group',
    component: GroupDetail
  },
  {
    path: '/transactions',
    name: 'transactions',
    component: () => import(/* webpackChunkName: "transactions" */ '../views/Transactions.vue')
  }
]

const routerOptions: RouterOptions = {
  mode: 'history',
  base: process.env.BASE_URL,
  routes: routes
}

export default new VueRouter(routerOptions)
