export type Page = 'dashboard' | 'subscriptions' | 'categories' | 'profile'

export interface User {
  id: number
  email: string
  name: string
}

export interface Category {
  id: number
  name: string
  description?: string | null
}

export interface Subscription {
  id: number
  category_id: number
  name: string
  description?: string | null
  cost: number
  currency?: string | null
  billing_cycle: string
  next_billing_date: string
  start_date: string
  end_date?: string | null
  status?: string | null
  auto_renew?: boolean | null
  payment_method?: string | null
  notes?: string | null
}

export interface SubForm {
  name: string
  category_id: number
  cost: number
  currency: string
  billing_cycle: string
  next_billing_date: string
  start_date: string
  status: string
  auto_renew: boolean
  description: string
  notes: string
  payment_method: string
}

export interface CatForm {
  name: string
  description: string
}
