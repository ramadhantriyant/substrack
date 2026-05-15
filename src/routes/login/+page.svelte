<script lang="ts">
	import { signInWithPopup, GoogleAuthProvider } from 'firebase/auth';
	import { auth } from '$lib/firebase';
	import { authState } from '$lib/stores/auth.svelte';
	import { goto } from '$app/navigation';

	let signingIn = $state(false);
	let error = $state('');

	$effect(() => {
		if (authState.isAuthed) void goto('/dashboard');
	});

	async function signIn() {
		signingIn = true;
		error = '';
		try {
			await signInWithPopup(auth, new GoogleAuthProvider());
			void goto('/dashboard');
		} catch (e: unknown) {
			const err = e as { code?: string };
			if (err.code !== 'auth/popup-closed-by-user') {
				error = 'Sign-in failed. Please try again.';
			}
		} finally {
			signingIn = false;
		}
	}
</script>

<div class="flex min-h-screen flex-col items-center justify-center bg-[#030712] px-4">
	<div class="w-full max-w-sm">
		<!-- Logo -->
		<div class="mb-10 text-center">
			<div class="mb-4 text-5xl font-bold tracking-tight text-[#f1f5f9]">
				<span class="text-[#f59e0b]">◎</span> substrack
			</div>
			<p class="text-[#64748b]">Track every subscription. Miss nothing.</p>
		</div>

		<!-- Card -->
		<div class="rounded-2xl border border-[#1e293b] bg-[#0a0f1e] p-8">
			<h1 class="mb-2 text-xl font-semibold text-[#f1f5f9]">Welcome back</h1>
			<p class="mb-6 text-sm text-[#64748b]">Sign in to manage your subscriptions</p>

			<button
				onclick={signIn}
				disabled={signingIn}
				class="flex w-full items-center justify-center gap-3 rounded-xl border border-[#1e293b] bg-[#0f172a] px-4 py-3 text-[#f1f5f9] transition-all hover:border-[#334155] hover:bg-[#1e293b] disabled:cursor-not-allowed disabled:opacity-60"
			>
				{#if signingIn}
					<div
						class="h-5 w-5 animate-spin rounded-full border-2 border-[#f59e0b] border-t-transparent"
					></div>
					<span>Signing in...</span>
				{:else}
					<svg class="h-5 w-5" viewBox="0 0 24 24">
						<path
							fill="#4285F4"
							d="M22.56 12.25c0-.78-.07-1.53-.2-2.25H12v4.26h5.92c-.26 1.37-1.04 2.53-2.21 3.31v2.77h3.57c2.08-1.92 3.28-4.74 3.28-8.09z"
						/>
						<path
							fill="#34A853"
							d="M12 23c2.97 0 5.46-.98 7.28-2.66l-3.57-2.77c-.98.66-2.23 1.06-3.71 1.06-2.86 0-5.29-1.93-6.16-4.53H2.18v2.84C3.99 20.53 7.7 23 12 23z"
						/>
						<path
							fill="#FBBC05"
							d="M5.84 14.09c-.22-.66-.35-1.36-.35-2.09s.13-1.43.35-2.09V7.07H2.18C1.43 8.55 1 10.22 1 12s.43 3.45 1.18 4.93l3.66-2.84z"
						/>
						<path
							fill="#EA4335"
							d="M12 5.38c1.62 0 3.06.56 4.21 1.64l3.15-3.15C17.45 2.09 14.97 1 12 1 7.7 1 3.99 3.47 2.18 7.07l3.66 2.84c.87-2.6 3.3-4.53 6.16-4.53z"
						/>
					</svg>
					<span class="font-medium">Continue with Google</span>
				{/if}
			</button>

			{#if error}
				<p class="mt-3 text-center text-sm text-red-400">{error}</p>
			{/if}
		</div>

		<p class="mt-6 text-center text-xs text-[#334155]">Your data is stored securely in Firebase</p>
	</div>
</div>
