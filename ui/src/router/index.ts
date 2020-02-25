import Vue from 'vue'
import VueRouter, { RouteConfig, RouterOptions } from 'vue-router'
import Home from '../views/Home.vue'
import Accounts from '@/views/Accounts.vue'
import AccountDetail from '@/views/Account.vue'
import GroupDetail from '@/views/Group.vue'
import AccountCreate from '@/views/Account/Create.vue'
import AccountCreateSingle from '@/views/Account/Create/Single.vue'
import AccountCreateMultiple from '@/views/Account/Create/Multiple.vue'
import AccountCreateUpload from '@/views/Account/Create/Upload.vue'
import GroupCreate from '@/views/Group/Create.vue'
import GroupCreateSingle from '@/views/Group/Create/Single.vue'
import GroupCreateMultiple from '@/views/Group/Create/Multiple.vue'
import GroupCreateUpload from '@/views/Group/Create/Upload.vue'

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
    path: '/accounts/create',
    component: AccountCreate,
    children: [{
      path: '',
      name: 'account_create',
      component: AccountCreateSingle
    }, {
      path: 'multiple',
      name: 'account_create_multiple',
      component: AccountCreateMultiple
    }, {
      path: 'upload',
      name: 'account_create_upload',
      component: AccountCreateUpload
    }]
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
    path: '/groups/create',
    component: GroupCreate,
    children: [{
      path: '',
      name: 'group_create',
      component: GroupCreateSingle
    }, {
      path: 'multiple',
      name: 'group_create_multiple',
      component: GroupCreateMultiple
    }, {
      path: 'upload',
      name: 'group_create_upload',
      component: GroupCreateUpload
    }]
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
