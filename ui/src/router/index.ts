import Vue from 'vue'
import VueRouter, { RouteConfig, RouterOptions } from 'vue-router'
import Home from '../views/Home.vue'
import Accounts from '@/views/Accounts.vue'
import Account from '@/views/Account.vue'
import AccountDetail from '@/views/Account/Detail.vue'
import TestView from '@/views/Test.vue'

Vue.use(VueRouter)

const routes: RouteConfig[] = [
  {
    path: '/accounts',
    name: 'accounts',
    component: Accounts
  },
  {
    path: '/account/:id',
    component: AccountDetail,
    children: [
      {
        path: 'edit',
        name: 'account_edit',
        component: TestView
      },
      {
        path: '',
        name: 'account',
        component: AccountDetail
      }]
  },
  {
    path: '/groups',
    name: 'groups',
    component: () => import(/* webpackChunkName: "groups" */ '../views/Groups.vue')
  },
  {
    path: '/group/:id',
    name: 'group',
    component: AccountDetail
  },
  {
    path: '/transactions',
    name: 'transactions',
    component: () => import(/* webpackChunkName: "transactions" */ '../views/Transactions.vue')
  },
  {
    path: '/',
    name: 'home',
    component: Home
  }
]

const routerOptions: RouterOptions = {
  mode: 'history',
  base: process.env.BASE_URL,
  routes: routes
}

export default new VueRouter(routerOptions)
