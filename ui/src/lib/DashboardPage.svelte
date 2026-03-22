<script lang="ts">
  import type { Subscription, Category } from './types.js'
  import { fmtCost, fmtDate, daysUntil, cycleLabel } from './helpers.js'

  let { subscriptions, categories, monthlyTotal, activeCount, upcomingBills, onNavigate, onAddSub }: {
    subscriptions: Subscription[]
    categories: Category[]
    monthlyTotal: number
    activeCount: number
    upcomingBills: Subscription[]
    onNavigate: (p: 'subscriptions') => void
    onAddSub: () => void
  } = $props()

  function catName(id: number) {
    return categories.find(c => c.id === id)?.name ?? '—'
  }
</script>

<div class="grid grid-cols-4 max-[900px]:grid-cols-2 gap-4 mb-7">
  <div class="bg-canvas border border-line rounded-[10px] py-[18px] px-5 shadow-card">
    <div class="text-xs text-fg uppercase tracking-[0.04em] mb-1.5">Monthly spend</div>
    <div class="text-[22px] font-semibold text-heading">{fmtCost(monthlyTotal, 'IDR')}</div>
  </div>
  <div class="bg-canvas border border-line rounded-[10px] py-[18px] px-5 shadow-card">
    <div class="text-xs text-fg uppercase tracking-[0.04em] mb-1.5">Annual spend</div>
    <div class="text-[22px] font-semibold text-heading">{fmtCost(monthlyTotal * 12, 'IDR')}</div>
  </div>
  <div class="bg-canvas border border-line rounded-[10px] py-[18px] px-5 shadow-card">
    <div class="text-xs text-fg uppercase tracking-[0.04em] mb-1.5">Active</div>
    <div class="text-[22px] font-semibold text-heading">{activeCount}</div>
  </div>
  <div class="bg-canvas border border-line rounded-[10px] py-[18px] px-5 shadow-card">
    <div class="text-xs text-fg uppercase tracking-[0.04em] mb-1.5">Due in 30 days</div>
    <div class="text-[22px] font-semibold text-heading">{upcomingBills.length}</div>
  </div>
</div>

<section class="bg-canvas border border-line rounded-[10px] p-[20px_24px] mb-5 shadow-card">
  <div class="flex items-center justify-between mb-4">
    <h2 class="text-[15px]">Upcoming renewals</h2>
  </div>
  {#if upcomingBills.length === 0}
    <p class="empty">No renewals in the next 30 days.</p>
  {:else}
    <div class="bg-canvas border border-line rounded-[10px] overflow-hidden shadow-card">
      {#each upcomingBills as sub (sub.id)}
        <div class="flex items-center justify-between px-5 py-[14px] border-b border-line gap-4 last:border-b-0">
          <div class="min-w-0">
            <div class="text-sm font-medium text-heading">{sub.name}</div>
            <div class="text-xs text-fg mt-0.5">{catName(sub.category_id)} · {cycleLabel(sub.billing_cycle)}</div>
          </div>
          <div class="flex items-center gap-3 shrink-0">
            <span class="text-sm font-semibold text-heading">{fmtCost(sub.cost, sub.currency ?? 'IDR')}</span>
            <span class="text-[11px] font-semibold px-2 py-[3px] rounded-full whitespace-nowrap {daysUntil(sub.next_billing_date) <= 3 ? 'bg-bad-bg text-bad' : 'bg-accent-bg text-accent'}">
              {daysUntil(sub.next_billing_date) === 0 ? 'Today' : `${daysUntil(sub.next_billing_date)}d`}
            </span>
          </div>
        </div>
      {/each}
    </div>
  {/if}
</section>

{#if subscriptions.length === 0}
  <div class="empty-state">
    <svg width="48" height="48" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5" opacity="0.3">
      <rect x="3" y="4" width="18" height="16" rx="2"/>
      <path d="M3 9h18M8 4v5M16 4v5"/>
    </svg>
    <p>No subscriptions yet.</p>
    <button class="btn-primary" onclick={onAddSub}>Add your first subscription</button>
  </div>
{/if}
