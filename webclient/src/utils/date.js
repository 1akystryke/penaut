export function parseCustomDate(dateString) {
  const date = new Date(dateString)
  if (isNaN(date.getTime())) throw new Error('Invalid date format')
  return date
}

export function formatMessengerDate(dateString) {
  const date = new Date(dateString)
  const now = new Date()
  if (isNaN(date.getTime())) throw new Error('Invalid date')

  const diffDays = Math.floor((now - date) / 1000 / 60 / 60 / 24)
  const hours = date.getHours().toString().padStart(2, '0')
  const minutes = date.getMinutes().toString().padStart(2, '0')

  if (diffDays < 1) return `${hours}:${minutes}`

  const day = date.getDate().toString().padStart(2, '0')
  const month = (date.getMonth() + 1).toString().padStart(2, '0')
  const year = date.getFullYear()

  if (year === now.getFullYear()) return `${day}.${month} в ${hours}:${minutes}`
  return `${day}.${month}.${year} в ${hours}:${minutes}`
}
