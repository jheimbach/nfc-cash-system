function formatMoney(value: number) {
  if (value) {
    return value.toFixed(2) + ' €'
  }
}

export default formatMoney
