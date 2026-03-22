<script lang="ts">
  import { untrack } from 'svelte'
  import type { Subscription, Category, SubForm } from './types.js'
  import { today } from './helpers.js'

  let { editingSub, categories, onSave, onClose }: {
    editingSub: Subscription | null
    categories: Category[]
    onSave: (form: SubForm) => Promise<void>
    onClose: () => void
  } = $props()

  let error = $state('')
  let loading = $state(false)

  let form = $state<SubForm>(untrack(() => editingSub ? {
    name: editingSub.name,
    category_id: editingSub.category_id,
    cost: editingSub.cost,
    currency: editingSub.currency ?? 'IDR',
    billing_cycle: editingSub.billing_cycle,
    next_billing_date: editingSub.next_billing_date.split('T')[0],
    start_date: editingSub.start_date.split('T')[0],
    status: editingSub.status ?? 'active',
    auto_renew: editingSub.auto_renew ?? true,
    description: editingSub.description ?? '',
    notes: editingSub.notes ?? '',
    payment_method: editingSub.payment_method ?? '',
  } : {
    name: '',
    category_id: categories[0]?.id ?? 0,
    cost: 0,
    currency: 'IDR',
    billing_cycle: 'monthly',
    next_billing_date: '',
    start_date: today(),
    status: 'active',
    auto_renew: true,
    description: '',
    notes: '',
    payment_method: '',
  }))

  async function handleSubmit() {
    loading = true; error = ''
    try {
      await onSave(form)
    } catch (e: unknown) {
      error = (e as Error).message
    } finally { loading = false }
  }
</script>

<div
  class="modal-overlay"
  role="presentation"
  onclick={onClose}
  onkeydown={(e) => { if (e.key === 'Escape') onClose() }}
>
  <div
    class="modal"
    role="dialog"
    aria-modal="true"
    tabindex="-1"
    onclick={(e) => e.stopPropagation()}
    onkeydown={(e) => e.stopPropagation()}
  >
    <div class="modal-header">
      <h2>{editingSub ? 'Edit subscription' : 'New subscription'}</h2>
      <button class="modal-close" onclick={onClose}>✕</button>
    </div>

    <form class="modal-body" onsubmit={(e) => { e.preventDefault(); handleSubmit() }}>
      <div class="form-grid">
        <label class="span-2">
          Name *
          <input bind:value={form.name} placeholder="Netflix, Spotify…" required />
        </label>

        <label>
          Category
          <select bind:value={form.category_id}>
            {#each categories as cat}
              <option value={cat.id}>{cat.name}</option>
            {/each}
            {#if categories.length === 0}
              <option disabled>No categories — add one first</option>
            {/if}
          </select>
        </label>

        <label>
          Status
          <select bind:value={form.status}>
            <option value="active">Active</option>
            <option value="paused">Paused</option>
            <option value="inactive">Inactive</option>
          </select>
        </label>

        <label>
          Cost *
          <input bind:value={form.cost} type="number" min="0" step="0.01" required />
        </label>

        <label>
          Currency
          <input bind:value={form.currency} placeholder="USD" maxlength="3" />
        </label>

        <label>
          Billing cycle
          <select bind:value={form.billing_cycle}>
            <option value="weekly">Weekly</option>
            <option value="monthly">Monthly</option>
            <option value="quarterly">Quarterly</option>
            <option value="biannual">Every 6 months</option>
            <option value="yearly">Yearly</option>
          </select>
        </label>

        <label>
          Next billing date *
          <input bind:value={form.next_billing_date} type="date" required />
        </label>

        <label>
          Start date
          <input bind:value={form.start_date} type="date" />
        </label>

        <label>
          Payment method
          <input bind:value={form.payment_method} placeholder="Credit card, PayPal…" />
        </label>

        <label class="span-2">
          Description
          <textarea bind:value={form.description} rows="2" placeholder="Optional notes about this subscription"></textarea>
        </label>

        <label class="checkbox span-2">
          <input type="checkbox" bind:checked={form.auto_renew} />
          Auto-renew
        </label>
      </div>

      {#if error}
        <div class="msg msg-error">{error}</div>
      {/if}

      <div class="modal-footer">
        <button type="button" class="btn-ghost" onclick={onClose}>Cancel</button>
        <button type="submit" class="btn-primary" disabled={loading}>
          {loading ? 'Saving…' : 'Save'}
        </button>
      </div>
    </form>
  </div>
</div>
