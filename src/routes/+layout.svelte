<script lang="ts">
	import './layout.css';
	import favicon from '$lib/assets/favicon.svg';
	import { authState } from '$lib/stores/auth.svelte';
	import { subsState } from '$lib/stores/subscriptions.svelte';
	import { goto } from '$app/navigation';
	import { page } from '$app/state';

	let { children } = $props();

	$effect(() => {
		const u = authState.user;
		if (u === null) return;
		if (u === false) {
			subsState.stopListener();
			// eslint-disable-next-line
			if (page.route.id !== '/login') void goto('/login');
		} else {
			subsState.startListener(u);
		}
	});
</script>

<svelte:head><link rel="icon" href={favicon} /></svelte:head>

{#if authState.user === null}
	<div class="flex h-screen items-center justify-center bg-[#030712]">
		<div class="h-8 w-8 animate-spin rounded-full border-2 border-[#f59e0b] border-t-transparent"></div>
	</div>
{:else}
	{@render children()}
{/if}
