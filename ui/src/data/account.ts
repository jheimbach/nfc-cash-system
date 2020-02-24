import Group from '@/data/group'
import { Named } from '@/data/globals'

export default interface Account extends Named {
  id: number,
  description?: string,
  saldo: number,
  nfcChipId: string,
  group: Group
}
