export function today() {
  return new Date().toISOString().split('T')[0]
}

export function toMonthly(cost: number, cycle: string) {
  switch (cycle) {
    case 'weekly':    return cost * 52 / 12
    case 'monthly':   return cost
    case 'quarterly': return cost / 3
    case 'biannual':  return cost / 6
    case 'yearly':    return cost / 12
    default:          return cost
  }
}

export function fmtCost(amount: number, currency = 'IDR') {
  try {
    return new Intl.NumberFormat('en-US', { style: 'currency', currency, minimumFractionDigits: 0, maximumFractionDigits: 0 }).format(amount)
  } catch {
    return `${currency} ${Math.round(amount)}`
  }
}

export function fmtDate(d: string) {
  return new Date(d).toLocaleDateString('en-US', { month: 'short', day: 'numeric', year: 'numeric' })
}

export function daysUntil(d: string) {
  return Math.ceil((new Date(d).getTime() - Date.now()) / 86_400_000)
}

export function cycleLabel(c: string) {
  return ({ weekly: 'Weekly', monthly: 'Monthly', quarterly: 'Quarterly', biannual: 'Every 6 mo', yearly: 'Yearly' } as Record<string, string>)[c] ?? c
}

export function statusColor(s?: string | null) {
  if (s === 'active') return 'green'
  if (s === 'paused') return 'yellow'
  return 'gray'
}
