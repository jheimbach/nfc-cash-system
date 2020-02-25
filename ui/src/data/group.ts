import { Named } from '@/data/globals'

export default interface Group extends Named {
  id: number,
  description?: string,
  canOverdraw?: boolean,
}

export const emptyGroup: Group = {
  name: '',
  id: 0,
  description: '',
  canOverdraw: false
}
