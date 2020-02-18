import Account from '@/data/account'

export interface Transaction {
  id: number,
  oldSaldo: number,
  newSaldo: number,
  amount: number,
  created: Date,
  account: Account
}
