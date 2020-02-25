import Group from '@/data/group'
import { Named } from '@/data/globals'

export default interface Account extends Named {
  id: number,
  description?: string,
  saldo: number,
  nfcChipId: string,
  group: Group | null
}
export const emptyAccount = {
  id: 0,
  name: '',
  description: '',
  group: null,
  saldo: 0,
  nfcChipId: ''
}
