<script lang="ts">
  import type { Category } from './types.js'

  let { categories, onAdd, onEdit, onDelete }: {
    categories: Category[]
    onAdd: () => void
    onEdit: (cat: Category) => void
    onDelete: (id: number) => Promise<void>
  } = $props()
</script>

<div class="flex items-center justify-between mb-5 gap-3">
  <span class="text-sm text-fg">{categories.length} {categories.length === 1 ? 'category' : 'categories'}</span>
  <button class="btn-primary" onclick={onAdd}>+ Add</button>
</div>

{#if categories.length === 0}
  <div class="empty-state">
    <p class="empty">No categories yet.</p>
    <button class="btn-primary" onclick={onAdd}>Add category</button>
  </div>
{:else}
  <div class="bg-canvas border border-line rounded-[10px] overflow-hidden shadow-card">
    {#each categories as cat (cat.id)}
      <div class="flex items-center justify-between px-5 py-[14px] border-b border-line gap-4 last:border-b-0">
        <div class="min-w-0">
          <div class="text-sm font-medium text-heading">{cat.name}</div>
          {#if cat.description}<div class="text-xs text-fg mt-0.5">{cat.description}</div>{/if}
        </div>
        <div class="flex items-center gap-2 shrink-0">
          <button class="btn-sm" onclick={() => onEdit(cat)}>Edit</button>
          <button class="btn-sm btn-danger" onclick={() => onDelete(cat.id)}>Delete</button>
        </div>
      </div>
    {/each}
  </div>
{/if}
