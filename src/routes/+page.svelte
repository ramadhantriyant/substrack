<script lang="ts">
	import { goto } from '$app/navigation';
	import { authState } from '$lib/stores/auth.svelte';

	$effect(() => {
		if (authState.user === null) return; // still loading
		void goto(authState.isAuthed ? '/dashboard' : '/login', { replaceState: true });
	});
</script>
