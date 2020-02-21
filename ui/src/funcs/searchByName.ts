import { Named } from '@/data/globals'

export function searchByName(items: Named[], term: string) {
  if (term) {
    return items.filter((item: Named) => {
      return toLower(item.name).includes(toLower(term))
    })
  }

  return items
}

function toLower(text: string) {
  return text.toString().toLowerCase()
}
