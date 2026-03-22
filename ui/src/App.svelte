<script lang="ts">
  import AuthPage from './lib/AuthPage.svelte'
  import Sidebar from './lib/Sidebar.svelte'
  import DashboardPage from './lib/DashboardPage.svelte'
  import SubscriptionsPage from './lib/SubscriptionsPage.svelte'
  import CategoriesPage from './lib/CategoriesPage.svelte'
  import ProfilePage from './lib/ProfilePage.svelte'
  import SubModal from './lib/SubModal.svelte'
  import CatModal from './lib/CatModal.svelte'
  import type { Page, User, Category, Subscription, SubForm, CatForm } from './lib/types.js'
  import { toMonthly, daysUntil } from './lib/helpers.js'

  // ── Auth ─────────────────────────────────────────────────────────────────────
  let token        = $state(localStorage.getItem('access_token') ?? '')
  let storedRefresh = $state(localStorage.getItem('refresh_token') ?? '')

  // ── App state ─────────────────────────────────────────────────────────────────
  let user          = $state<User | null>(null)
  let subscriptions = $state<Subscription[]>([])
  let categories    = $state<Category[]>([])
  let page          = $state<Page>('dashboard')
  let globalError   = $state('')
  let sidebarOpen   = $state(false)

  // ── Modals ────────────────────────────────────────────────────────────────────
  let showSubModal = $state(false)
  let editingSub   = $state<Subscription | null>(null)
  let showCatModal = $state(false)
  let editingCat   = $state<Category | null>(null)

  // ── Derived ───────────────────────────────────────────────────────────────────
  const isAuthed = $derived(!!token)

  const monthlyTotal = $derived(
    subscriptions
      .filter(s => s.status === 'active')
      .reduce((sum, s) => sum + toMonthly(s.cost, s.billing_cycle), 0)
  )

  const activeCount = $derived(subscriptions.filter(s => s.status === 'active').length)

  const upcomingBills = $derived(
    subscriptions
      .filter(s => s.status === 'active' && daysUntil(s.next_billing_date) >= 0 && daysUntil(s.next_billing_date) <= 30)
      .sort((a, b) => new Date(a.next_billing_date).getTime() - new Date(b.next_billing_date).getTime())
  )

  // ── API ───────────────────────────────────────────────────────────────────────
  async function api<T = unknown>(path: string, opts: RequestInit = {}): Promise<T | null> {
    const res = await fetch(path, {
      ...opts,
      headers: {
        'Content-Type': 'application/json',
        ...(token ? { Authorization: `Bearer ${token}` } : {}),
        ...(opts.headers as Record<string, string> ?? {}),
      },
    })
    if (res.status === 401) { doLogout(); return null }
    if (!res.ok) {
      const body = await res.json().catch(() => ({}))
      throw new Error(body.error ?? body.message ?? res.statusText)
    }
    return res.json()
  }

  // ── Bootstrap ─────────────────────────────────────────────────────────────────
  $effect(() => { if (token) loadAll() })

  async function loadAll() {
    await Promise.all([loadUser(), loadSubscriptions(), loadCategories()])
  }

  async function loadUser() {
    try { user = await api('/api/user/me') } catch {}
  }

  async function loadSubscriptions() {
    try {
      const data = await api<{ subscriptions: Subscription[] }>('/api/user/me/subscription')
      subscriptions = data?.subscriptions ?? []
    } catch {}
  }

  async function loadCategories() {
    try {
      const data = await api<{ categories: Category[] }>('/api/user/me/category')
      categories = data?.categories ?? []
    } catch {}
  }

  // ── Auth actions ──────────────────────────────────────────────────────────────
  function setTokens(access: string, refresh: string) {
    token = access
    storedRefresh = refresh
    localStorage.setItem('access_token', access)
    localStorage.setItem('refresh_token', refresh)
  }

  async function doLogout() {
    try {
      await api('/auth/logout', { method: 'POST', body: JSON.stringify({ refresh_token: storedRefresh }) })
    } catch {}
    token = ''; storedRefresh = ''
    user = null; subscriptions = []; categories = []
    localStorage.removeItem('access_token')
    localStorage.removeItem('refresh_token')
    page = 'dashboard'
  }

  // ── Subscription actions ──────────────────────────────────────────────────────
  function openSubModal(sub?: Subscription) {
    editingSub = sub ?? null
    showSubModal = true
  }

  async function saveSub(form: SubForm) {
    const body = {
      ...form,
      cost: Number(form.cost),
      category_id: Number(form.category_id),
      next_billing_date: new Date(form.next_billing_date).toISOString(),
      start_date: new Date(form.start_date).toISOString(),
      description: form.description || null,
      notes: form.notes || null,
      payment_method: form.payment_method || null,
    }
    if (editingSub) {
      await api(`/api/subscription/${editingSub.id}`, { method: 'PUT', body: JSON.stringify(body) })
    } else {
      const created = await api<Subscription>('/api/subscription', { method: 'POST', body: JSON.stringify(body) })
      if (created) await api(`/api/user/me/subscription/${created.id}`, { method: 'POST' })
    }
    await loadSubscriptions()
    showSubModal = false
    editingSub = null
  }

  async function deleteSub(id: number) {
    if (!confirm('Delete this subscription?')) return
    try {
      await api(`/api/user/me/subscription/${id}`, { method: 'DELETE' })
      await loadSubscriptions()
    } catch (e: unknown) { globalError = (e as Error).message }
  }

  // ── Category actions ──────────────────────────────────────────────────────────
  function openCatModal(cat?: Category) {
    editingCat = cat ?? null
    showCatModal = true
  }

  async function saveCat(form: CatForm) {
    if (editingCat) {
      await api(`/api/category/${editingCat.id}`, { method: 'PUT', body: JSON.stringify(form) })
    } else {
      let cat: Category | null = null
      try {
        cat = await api<Category>('/api/category', { method: 'POST', body: JSON.stringify(form) })
      } catch (e: unknown) {
        if ((e as Error).message !== 'category name already exists') throw e
        cat = await api<Category>(`/api/category/name/${encodeURIComponent(form.name)}`)
      }
      if (cat) {
        try {
          await api(`/api/user/me/category/${cat.id}`, { method: 'POST' })
        } catch (e: unknown) {
          if ((e as Error).message !== 'category already added') throw e
        }
      }
    }
    await loadCategories()
    showCatModal = false
    editingCat = null
  }

  async function deleteCat(id: number) {
    if (!confirm('Delete this category?')) return
    try {
      await api(`/api/user/me/category/${id}`, { method: 'DELETE' })
      await loadCategories()
    } catch (e: unknown) { globalError = (e as Error).message }
  }

  // ── Profile actions ───────────────────────────────────────────────────────────
  async function updateProfile(form: { email: string; name: string }) {
    const updated = await api<User>('/api/user/me', { method: 'PUT', body: JSON.stringify(form) })
    if (updated) user = updated
  }

  async function updatePassword(form: { old_password: string; new_password: string }) {
    await api('/api/user/me/password', { method: 'PUT', body: JSON.stringify(form) })
  }

  // ── Navigation ────────────────────────────────────────────────────────────────
  function navigate(p: Page) {
    page = p
    sidebarOpen = false
    globalError = ''
  }
</script>

{#if !isAuthed}
  <AuthPage onLogin={setTokens} />
{:else}
  {#if sidebarOpen}
    <div
      class="fixed inset-0 bg-black/40 z-40 hidden max-md:block"
      role="presentation"
      onclick={() => sidebarOpen = false}
      onkeydown={(e) => { if (e.key === 'Escape') sidebarOpen = false }}
    ></div>
  {/if}

  <div class="flex min-h-svh">
    <Sidebar {page} {user} open={sidebarOpen} onNavigate={navigate} onLogout={doLogout} />

    <div class="flex-1 min-w-0 flex flex-col">
      <header class="flex items-center gap-3 px-6 h-14 border-b border-line bg-canvas sticky top-0 z-10">
        <button
          class="hidden max-md:flex items-center justify-center w-9 h-9 border-0 rounded-[7px] bg-transparent text-fg hover:bg-accent-bg"
          onclick={() => sidebarOpen = !sidebarOpen}
          aria-label="Menu"
        >
          <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <path d="M3 6h18M3 12h18M3 18h18"/>
          </svg>
        </button>
        <span class="font-semibold text-[15px] text-heading flex-1">
          {#if page === 'dashboard'}Dashboard
          {:else if page === 'subscriptions'}Subscriptions
          {:else if page === 'categories'}Categories
          {:else}Profile{/if}
        </span>
        <span class="text-[13px] text-fg">{user?.name ?? ''}</span>
      </header>

      {#if globalError}
        <div class="msg msg-error msg-bar">
          {globalError}
          <button onclick={() => globalError = ''}>✕</button>
        </div>
      {/if}

      <div class="p-6 flex-1">
        {#if page === 'dashboard'}
          <DashboardPage
            {subscriptions}
            {categories}
            {monthlyTotal}
            {activeCount}
            {upcomingBills}
            onNavigate={navigate}
            onAddSub={() => openSubModal()}
          />
        {:else if page === 'subscriptions'}
          <SubscriptionsPage
            {subscriptions}
            {categories}
            onAdd={() => openSubModal()}
            onEdit={openSubModal}
            onDelete={deleteSub}
          />
        {:else if page === 'categories'}
          <CategoriesPage
            {categories}
            onAdd={() => openCatModal()}
            onEdit={openCatModal}
            onDelete={deleteCat}
          />
        {:else if page === 'profile'}
          <ProfilePage
            {user}
            onUpdateProfile={updateProfile}
            onUpdatePassword={updatePassword}
          />
        {/if}
      </div>
    </div>
  </div>

  {#if showSubModal}
    <SubModal
      {editingSub}
      {categories}
      onSave={saveSub}
      onClose={() => { showSubModal = false; editingSub = null }}
    />
  {/if}

  {#if showCatModal}
    <CatModal
      {editingCat}
      onSave={saveCat}
      onClose={() => { showCatModal = false; editingCat = null }}
    />
  {/if}
{/if}
