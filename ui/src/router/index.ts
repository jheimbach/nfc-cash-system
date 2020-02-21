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
    component: Account,
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
  }
]

const routerOptions: RouterOptions = {
  mode: 'history',
  base: process.env.BASE_URL,
  routes: routes
}

const router = (linkClass: string) => {
  routerOptions.linkActiveClass = linkClass
  return new VueRouter(routerOptions)
}

export default router
