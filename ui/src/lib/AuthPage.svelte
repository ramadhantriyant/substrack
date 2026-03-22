<script lang="ts">
  let { onLogin }: { onLogin: (access: string, refresh: string) => void } = $props()

  let mode = $state<'login' | 'register'>('login')
  let error = $state('')
  let loading = $state(false)
  let loginForm = $state({ email: '', password: '' })
  let registerForm = $state({ email: '', name: '', password: '' })

  async function fetchJson(path: string, body: unknown) {
    const res = await fetch(path, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(body),
    })
    const data = await res.json().catch(() => ({}))
    if (!res.ok) throw new Error(data.error ?? data.message ?? res.statusText)
    return data
  }

  async function doLogin() {
    loading = true; error = ''
    try {
      const data = await fetchJson('/auth/login', loginForm)
      onLogin(data.access_token, data.refresh_token)
    } catch (e: unknown) {
      error = (e as Error).message
    } finally { loading = false }
  }

  async function doRegister() {
    loading = true; error = ''
    try {
      await fetchJson('/auth/register', registerForm)
      loginForm = { email: registerForm.email, password: registerForm.password }
      await doLogin()
    } catch (e: unknown) {
      error = (e as Error).message
    } finally { loading = false }
  }
</script>

<div class="min-h-svh flex items-center justify-center p-5 bg-surface">
  <div class="w-full max-w-[400px] bg-canvas border border-line rounded-xl p-9 shadow-popup">
    <div class="text-center mb-7">
      <div class="inline-flex items-center justify-center w-[52px] h-[52px] bg-accent-bg border border-accent-ring rounded-xl text-accent mb-3">
        <svg width="28" height="28" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
          <rect x="3" y="4" width="18" height="16" rx="2"/>
          <path d="M3 9h18M8 4v5M16 4v5"/>
        </svg>
      </div>
      <h1 class="text-[22px] mb-1">Substrack</h1>
      <p class="text-sm text-fg">Track every subscription you own</p>
    </div>

    <div class="flex border border-line rounded-lg overflow-hidden mb-5">
      <button
        class="flex-1 py-2 text-sm font-medium border-0 transition-all duration-150 cursor-pointer {mode === 'login' ? 'bg-accent text-white' : 'bg-transparent text-fg'}"
        onclick={() => { mode = 'login'; error = '' }}
      >Sign in</button>
      <button
        class="flex-1 py-2 text-sm font-medium border-0 transition-all duration-150 cursor-pointer {mode === 'register' ? 'bg-accent text-white' : 'bg-transparent text-fg'}"
        onclick={() => { mode = 'register'; error = '' }}
      >Create account</button>
    </div>

    {#if error}
      <div class="msg msg-error">{error}</div>
    {/if}

    {#if mode === 'login'}
      <form class="flex flex-col gap-3.5" onsubmit={(e) => { e.preventDefault(); doLogin() }}>
        <label>
          Email
          <input bind:value={loginForm.email} type="email" placeholder="you@example.com" required autocomplete="email" />
        </label>
        <label>
          Password
          <input bind:value={loginForm.password} type="password" placeholder="••••••••" required autocomplete="current-password" />
        </label>
        <button type="submit" class="btn-primary" disabled={loading}>
          {loading ? 'Signing in…' : 'Sign in'}
        </button>
      </form>
    {:else}
      <form class="flex flex-col gap-3.5" onsubmit={(e) => { e.preventDefault(); doRegister() }}>
        <label>
          Full name
          <input bind:value={registerForm.name} placeholder="Your name" required autocomplete="name" />
        </label>
        <label>
          Email
          <input bind:value={registerForm.email} type="email" placeholder="you@example.com" required autocomplete="email" />
        </label>
        <label>
          Password
          <input bind:value={registerForm.password} type="password" placeholder="••••••••" required autocomplete="new-password" />
        </label>
        <button type="submit" class="btn-primary" disabled={loading}>
          {loading ? 'Creating account…' : 'Create account'}
        </button>
      </form>
    {/if}
  </div>
</div>
