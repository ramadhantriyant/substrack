<script lang="ts">
  import { untrack } from 'svelte'
  import type { Category, CatForm } from './types.js'

  let { editingCat, onSave, onClose }: {
    editingCat: Category | null
    onSave: (form: CatForm) => Promise<void>
    onClose: () => void
  } = $props()

  let error = $state('')
  let loading = $state(false)
  let form = $state<CatForm>(untrack(() => ({ name: editingCat?.name ?? '', description: editingCat?.description ?? '' })))

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
    class="modal modal-sm"
    role="dialog"
    aria-modal="true"
    tabindex="-1"
    onclick={(e) => e.stopPropagation()}
    onkeydown={(e) => e.stopPropagation()}
  >
    <div class="modal-header">
      <h2>{editingCat ? 'Edit category' : 'New category'}</h2>
      <button class="modal-close" onclick={onClose}>✕</button>
    </div>

    <form class="modal-body" onsubmit={(e) => { e.preventDefault(); handleSubmit() }}>
      <label>
        Name *
        <input bind:value={form.name} required />
      </label>
      <label>
        Description
        <input bind:value={form.description} placeholder="Optional" />
      </label>

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
