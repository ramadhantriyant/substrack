<script lang="ts">
	import { goto } from '$app/navigation';

	$effect(() => {
		// eslint-disable-next-line
		void goto('/dashboard');
	});
</script>
