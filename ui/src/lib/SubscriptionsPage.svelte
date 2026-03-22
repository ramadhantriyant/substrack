<script lang="ts">
  import type { Subscription, Category } from './types.js'
  import { fmtCost, fmtDate, cycleLabel, statusColor } from './helpers.js'

  let { subscriptions, categories, onAdd, onEdit, onDelete }: {
    subscriptions: Subscription[]
    categories: Category[]
    onAdd: () => void
    onEdit: (sub: Subscription) => void
    onDelete: (id: number) => Promise<void>
  } = $props()

  let filter = $state<'all' | 'active' | 'inactive'>('all')

  const filtered = $derived(
    filter === 'active'   ? subscriptions.filter(s => s.status === 'active') :
    filter === 'inactive' ? subscriptions.filter(s => s.status !== 'active') :
    subscriptions
  )

  function catName(id: number) {
    return categories.find(c => c.id === id)?.name ?? '—'
  }
</script>

<div class="flex items-center justify-between mb-5 gap-3 flex-wrap">
  <div class="flex gap-1 bg-surface border border-line rounded-lg p-1">
    <button
      class="py-[5px] px-3 text-[13px] border-0 rounded-[5px] transition-all duration-150 whitespace-nowrap cursor-pointer {filter === 'all' ? 'bg-canvas text-heading font-medium shadow-card' : 'bg-transparent text-fg'}"
      onclick={() => filter = 'all'}
    >All ({subscriptions.length})</button>
    <button
      class="py-[5px] px-3 text-[13px] border-0 rounded-[5px] transition-all duration-150 whitespace-nowrap cursor-pointer {filter === 'active' ? 'bg-canvas text-heading font-medium shadow-card' : 'bg-transparent text-fg'}"
      onclick={() => filter = 'active'}
    >Active ({subscriptions.filter(s => s.status === 'active').length})</button>
    <button
      class="py-[5px] px-3 text-[13px] border-0 rounded-[5px] transition-all duration-150 whitespace-nowrap cursor-pointer {filter === 'inactive' ? 'bg-canvas text-heading font-medium shadow-card' : 'bg-transparent text-fg'}"
      onclick={() => filter = 'inactive'}
    >Inactive ({subscriptions.filter(s => s.status !== 'active').length})</button>
  </div>
  <button class="btn-primary" onclick={onAdd}>+ Add</button>
</div>

{#if filtered.length === 0}
  <div class="empty-state">
    <p class="empty">No subscriptions here.</p>
    {#if filter === 'all'}
      <button class="btn-primary" onclick={onAdd}>Add subscription</button>
    {/if}
  </div>
{:else}
  <div class="grid grid-cols-[repeat(auto-fill,minmax(280px,1fr))] gap-4">
    {#each filtered as sub (sub.id)}
      <div class="bg-canvas border border-line rounded-[10px] py-[18px] px-5 flex flex-col gap-[10px] shadow-card hover:shadow-popup transition-shadow duration-200">
        <div class="flex items-start justify-between gap-[10px]">
          <div>
            <div class="text-[15px] font-semibold text-heading">{sub.name}</div>
            <div class="text-xs text-fg mt-0.5">{catName(sub.category_id)}</div>
          </div>
          <span class="badge badge-{statusColor(sub.status)}">{sub.status ?? 'inactive'}</span>
        </div>
        <div class="text-xl font-semibold text-heading">
          {fmtCost(sub.cost, sub.currency ?? 'IDR')}
          <span class="text-[13px] font-normal text-fg">/ {cycleLabel(sub.billing_cycle).toLowerCase()}</span>
        </div>
        {#if sub.description}
          <div class="text-[13px] text-fg leading-snug">{sub.description}</div>
        {/if}
        <div class="flex items-center justify-between mt-1 gap-2">
          <div class="flex items-center gap-[5px] text-xs text-fg">
            <svg width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <rect x="3" y="4" width="18" height="16" rx="2"/><path d="M3 9h18M8 4v5M16 4v5"/>
            </svg>
            {fmtDate(sub.next_billing_date)}
          </div>
          <div class="flex gap-1.5">
            <button class="btn-sm" onclick={() => onEdit(sub)}>Edit</button>
            <button class="btn-sm btn-danger" onclick={() => onDelete(sub.id)}>Delete</button>
          </div>
        </div>
      </div>
    {/each}
  </div>
{/if}
