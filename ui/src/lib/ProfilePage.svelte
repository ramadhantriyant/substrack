<script lang="ts">
  import { untrack } from 'svelte'
  import type { User } from './types.js'

  let { user, onUpdateProfile, onUpdatePassword }: {
    user: User | null
    onUpdateProfile: (form: { email: string; name: string }) => Promise<void>
    onUpdatePassword: (form: { old_password: string; new_password: string }) => Promise<void>
  } = $props()

  let profileForm = $state(untrack(() => ({ email: user?.email ?? '', name: user?.name ?? '' })))
  let passwordForm = $state({ old_password: '', new_password: '', confirm: '' })
  let profileError  = $state('')
  let passwordError = $state('')
  let profileLoading  = $state(false)
  let passwordLoading = $state(false)

  $effect(() => {
    if (user) profileForm = { email: user.email, name: user.name }
  })

  async function handleProfile() {
    profileLoading = true; profileError = ''
    try {
      await onUpdateProfile(profileForm)
    } catch (e: unknown) {
      profileError = (e as Error).message
    } finally { profileLoading = false }
  }

  async function handlePassword() {
    if (passwordForm.new_password !== passwordForm.confirm) {
      passwordError = 'Passwords do not match'; return
    }
    passwordLoading = true; passwordError = ''
    try {
      await onUpdatePassword({ old_password: passwordForm.old_password, new_password: passwordForm.new_password })
      passwordForm = { old_password: '', new_password: '', confirm: '' }
    } catch (e: unknown) {
      passwordError = (e as Error).message
    } finally { passwordLoading = false }
  }
</script>

<div class="grid grid-cols-[repeat(auto-fill,minmax(320px,1fr))] gap-5 items-start">
  <section class="bg-canvas border border-line rounded-[10px] p-[20px_24px] shadow-card">
    <h2 class="text-[15px] mb-4">Account</h2>
    <form class="flex flex-col gap-3.5" onsubmit={(e) => { e.preventDefault(); handleProfile() }}>
      <label>
        Name
        <input bind:value={profileForm.name} required />
      </label>
      <label>
        Email
        <input bind:value={profileForm.email} type="email" required />
      </label>
      {#if profileError}
        <div class="msg msg-error">{profileError}</div>
      {/if}
      <button type="submit" class="btn-primary" disabled={profileLoading}>
        {profileLoading ? 'Saving…' : 'Save changes'}
      </button>
    </form>
  </section>

  <section class="bg-canvas border border-line rounded-[10px] p-[20px_24px] shadow-card">
    <h2 class="text-[15px] mb-4">Change password</h2>
    <form class="flex flex-col gap-3.5" onsubmit={(e) => { e.preventDefault(); handlePassword() }}>
      <label>
        Current password
        <input bind:value={passwordForm.old_password} type="password" required autocomplete="current-password" />
      </label>
      <label>
        New password
        <input bind:value={passwordForm.new_password} type="password" required autocomplete="new-password" />
      </label>
      <label>
        Confirm new password
        <input bind:value={passwordForm.confirm} type="password" required autocomplete="new-password" />
      </label>
      {#if passwordError}
        <div class="msg msg-error">{passwordError}</div>
      {/if}
      <button type="submit" class="btn-primary" disabled={passwordLoading}>
        {passwordLoading ? 'Updating…' : 'Update password'}
      </button>
    </form>
  </section>
</div>
