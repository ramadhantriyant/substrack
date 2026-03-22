<script lang="ts">
  import type { Page, User } from './types.js'

  let { page, user, open, onNavigate, onLogout }: {
    page: Page
    user: User | null
    open: boolean
    onNavigate: (p: Page) => void
    onLogout: () => void
  } = $props()
</script>

<nav class="w-[220px] shrink-0 bg-sidebar border-r border-line flex flex-col sticky top-0 h-svh max-md:fixed max-md:left-0 max-md:top-0 max-md:bottom-0 max-md:z-50 max-md:transition-transform max-md:duration-250 {open ? 'max-md:translate-x-0' : 'max-md:-translate-x-full'}">
  <div class="flex items-center gap-[10px] px-4 py-[18px] border-b border-line font-semibold text-[15px] text-heading">
    <div class="flex items-center justify-center w-[30px] h-[30px] bg-accent-bg rounded-[7px] text-accent shrink-0">
      <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5">
        <rect x="3" y="4" width="18" height="16" rx="2"/>
        <path d="M3 9h18M8 4v5M16 4v5"/>
      </svg>
    </div>
    <span>Substrack</span>
  </div>

  <div class="flex-1 p-[10px_8px] flex flex-col gap-0.5">
    <button
      class="flex items-center gap-[10px] py-[9px] px-[10px] border-0 rounded-[7px] text-sm text-left w-full transition-all duration-150 cursor-pointer {page === 'dashboard' ? 'bg-accent-bg text-accent font-medium' : 'bg-transparent text-fg hover:bg-accent-bg hover:text-heading'}"
      onclick={() => onNavigate('dashboard')}
    >
      <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
        <rect x="3" y="3" width="7" height="7"/><rect x="14" y="3" width="7" height="7"/>
        <rect x="3" y="14" width="7" height="7"/><rect x="14" y="14" width="7" height="7"/>
      </svg>
      Dashboard
    </button>
    <button
      class="flex items-center gap-[10px] py-[9px] px-[10px] border-0 rounded-[7px] text-sm text-left w-full transition-all duration-150 cursor-pointer {page === 'subscriptions' ? 'bg-accent-bg text-accent font-medium' : 'bg-transparent text-fg hover:bg-accent-bg hover:text-heading'}"
      onclick={() => onNavigate('subscriptions')}
    >
      <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
        <path d="M12 2C6.48 2 2 6.48 2 12s4.48 10 10 10 10-4.48 10-10S17.52 2 12 2z"/>
        <path d="M12 6v6l4 2"/>
      </svg>
      Subscriptions
    </button>
    <button
      class="flex items-center gap-[10px] py-[9px] px-[10px] border-0 rounded-[7px] text-sm text-left w-full transition-all duration-150 cursor-pointer {page === 'categories' ? 'bg-accent-bg text-accent font-medium' : 'bg-transparent text-fg hover:bg-accent-bg hover:text-heading'}"
      onclick={() => onNavigate('categories')}
    >
      <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
        <path d="M4 6h16M4 12h16M4 18h7"/>
      </svg>
      Categories
    </button>
    <button
      class="flex items-center gap-[10px] py-[9px] px-[10px] border-0 rounded-[7px] text-sm text-left w-full transition-all duration-150 cursor-pointer {page === 'profile' ? 'bg-accent-bg text-accent font-medium' : 'bg-transparent text-fg hover:bg-accent-bg hover:text-heading'}"
      onclick={() => onNavigate('profile')}
    >
      <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
        <circle cx="12" cy="8" r="4"/>
        <path d="M4 20c0-4 3.6-7 8-7s8 3 8 7"/>
      </svg>
      Profile
    </button>
  </div>

  <button
    class="flex items-center gap-[10px] py-3 px-[18px] border-0 border-t border-line bg-transparent text-fg text-[13px] cursor-pointer transition-colors duration-150 w-full text-left whitespace-nowrap overflow-hidden text-ellipsis hover:text-bad"
    onclick={onLogout}
  >
    <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
      <path d="M9 21H5a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h4M16 17l5-5-5-5M21 12H9"/>
    </svg>
    Sign out · {user?.name ?? ''}
  </button>
</nav>
