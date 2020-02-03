import moment from 'moment'

function formatDate (value: string) {
  if (value) {
    return moment(value).format('YYYY-MM-DD hh:mm')
  }
}

export default formatDate
