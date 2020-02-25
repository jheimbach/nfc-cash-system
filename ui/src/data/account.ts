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
  name: '11',
  description: '',
  group: { id: 1, name: 'test' },
  saldo: 0,
  nfcChipId: '11'
}
