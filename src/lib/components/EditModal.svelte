<script lang="ts">
	import { subsState } from '$lib/stores/subscriptions.svelte';
	import type { Subscription } from '$lib/types';

	let {
		onClose,
		uid,
		subscription
	}: { onClose: () => void; uid: string; subscription: Subscription } = $props();

	const ICONS = [
		'📱',
		'🎬',
		'🎵',
		'📺',
		'🎮',
		'💻',
		'☁️',
		'📧',
		'🔒',
		'💪',
		'🏥',
		'💳',
		'📊',
		'🛡️',
		'🌐',
		'📰',
		'🎨',
		'📚',
		'🎯',
		'⚡',
		'🎙️',
		'🎧',
		'✈️',
		'🏠'
	];
	const COLORS = [
		'#f59e0b',
		'#6366f1',
		'#10b981',
		'#ec4899',
		'#3b82f6',
		'#8b5cf6',
		'#ef4444',
		'#14b8a6',
		'#f97316',
		'#84cc16',
		'#06b6d4',
		'#d946ef'
	];

	let name = $state(subscription.name);
	let amount = $state(String(subscription.amount));
	let cycle = $state<Subscription['cycle']>(subscription.cycle);
	let category = $state<Subscription['category']>(subscription.category);
	let nextBilling = $state(subscription.nextBilling);
	let icon = $state(subscription.icon);
	let color = $state(subscription.color);
	let error = $state('');
	let submitting = $state(false);

	async function handleSubmit() {
		if (!name.trim()) {
			error = 'Name is required';
			return;
		}
		if (!amount || Number(amount) <= 0) {
			error = 'Amount must be greater than 0';
			return;
		}
		if (!nextBilling) {
			error = 'Next billing date is required';
			return;
		}
		error = '';
		submitting = true;
		try {
			await subsState.update(uid, subscription.id, {
				name: name.trim(),
				amount: Number(amount),
				cycle,
				category,
				nextBilling,
				icon,
				color
			});
			onClose();
		} catch {
			error = 'Failed to update subscription. Please try again.';
		} finally {
			submitting = false;
		}
	}

	function handleBackdrop(e: MouseEvent) {
		if (e.target === e.currentTarget) onClose();
	}

	function handleKeydown(e: KeyboardEvent) {
		if (e.key === 'Escape') onClose();
	}
</script>

<div
	class="fixed inset-0 z-50 flex items-center justify-center bg-black/60 p-4 backdrop-blur-sm"
	role="dialog"
	aria-modal="true"
	aria-label="Edit Subscription"
	tabindex="-1"
	onclick={handleBackdrop}
	onkeydown={handleKeydown}
>
	<div class="w-full max-w-md rounded-2xl border border-[#1e293b] bg-[#0a0f1e] p-6 shadow-2xl">
		<div class="mb-6 flex items-center justify-between">
			<h2 class="text-lg font-semibold text-[#f1f5f9]">Edit Subscription</h2>
			<button
				onclick={onClose}
				aria-label="Close"
				class="text-xl leading-none text-[#64748b] transition-colors hover:text-[#f1f5f9]">×</button
			>
		</div>

		<form
			onsubmit={(e) => {
				e.preventDefault();
				handleSubmit();
			}}
			class="space-y-4"
		>
			<div>
				<label for="edit-name" class="mb-1 block text-sm text-[#64748b]">Service Name</label>
				<input
					id="edit-name"
					type="text"
					bind:value={name}
					placeholder="Netflix, Spotify..."
					class="w-full rounded-lg border border-[#1e293b] bg-[#0f172a] px-3 py-2 text-[#f1f5f9] transition-colors placeholder:text-[#334155] focus:border-[#f59e0b] focus:outline-none"
				/>
			</div>

			<div class="grid grid-cols-2 gap-3">
				<div>
					<label for="edit-amount" class="mb-1 block text-sm text-[#64748b]">Amount (IDR)</label>
					<input
						id="edit-amount"
						type="number"
						bind:value={amount}
						placeholder="50000"
						min="0"
						class="w-full rounded-lg border border-[#1e293b] bg-[#0f172a] px-3 py-2 font-mono text-[#f1f5f9] transition-colors placeholder:text-[#334155] focus:border-[#f59e0b] focus:outline-none"
					/>
				</div>
				<div>
					<label for="edit-cycle" class="mb-1 block text-sm text-[#64748b]">Cycle</label>
					<select
						id="edit-cycle"
						bind:value={cycle}
						class="w-full rounded-lg border border-[#1e293b] bg-[#0f172a] px-3 py-2 text-[#f1f5f9] transition-colors focus:border-[#f59e0b] focus:outline-none"
					>
						<option value="Monthly">Monthly</option>
						<option value="Yearly">Yearly</option>
						<option value="Weekly">Weekly</option>
					</select>
				</div>
			</div>

			<div>
				<label for="edit-category" class="mb-1 block text-sm text-[#64748b]">Category</label>
				<select
					id="edit-category"
					bind:value={category}
					class="w-full rounded-lg border border-[#1e293b] bg-[#0f172a] px-3 py-2 text-[#f1f5f9] transition-colors focus:border-[#f59e0b] focus:outline-none"
				>
					<option value="Entertainment">Entertainment</option>
					<option value="SaaS">SaaS</option>
					<option value="Utilities">Utilities</option>
					<option value="Health">Health</option>
					<option value="Finance">Finance</option>
					<option value="Other">Other</option>
				</select>
			</div>

			<div>
				<label for="edit-billing" class="mb-1 block text-sm text-[#64748b]"
					>Next Billing Date</label
				>
				<input
					id="edit-billing"
					type="date"
					bind:value={nextBilling}
					class="w-full rounded-lg border border-[#1e293b] bg-[#0f172a] px-3 py-2 text-[#f1f5f9] scheme-dark transition-colors focus:border-[#f59e0b] focus:outline-none"
				/>
			</div>

			<div>
				<p class="mb-2 text-sm text-[#64748b]">Icon</p>
				<div class="grid grid-cols-8 gap-1">
					{#each ICONS as ic (ic)}
						<button
							type="button"
							onclick={() => (icon = ic)}
							aria-label="Select icon {ic}"
							aria-pressed={icon === ic}
							class="flex h-8 w-8 items-center justify-center rounded-lg text-base transition-colors {icon ===
							ic
								? 'bg-[#f59e0b]/20 ring-1 ring-[#f59e0b]'
								: 'hover:bg-[#1e293b]'}">{ic}</button
						>
					{/each}
				</div>
			</div>

			<div>
				<p class="mb-2 text-sm text-[#64748b]">Color</p>
				<div class="flex flex-wrap gap-2">
					{#each COLORS as c (c)}
						<button
							type="button"
							onclick={() => (color = c)}
							aria-label="Select color {c}"
							aria-pressed={color === c}
							class="h-6 w-6 rounded-full transition-transform {color === c
								? 'scale-125 ring-2 ring-white/50'
								: 'hover:scale-110'}"
							style="background-color: {c}"
						></button>
					{/each}
				</div>
			</div>

			{#if error}
				<p class="text-sm text-red-400" role="alert">{error}</p>
			{/if}

			<div class="flex gap-3 pt-2">
				<button
					type="button"
					onclick={onClose}
					class="flex-1 rounded-lg border border-[#1e293b] py-2 text-[#64748b] transition-colors hover:border-[#334155] hover:text-[#f1f5f9]"
					>Cancel</button
				>
				<button
					type="submit"
					disabled={submitting}
					class="flex-1 rounded-lg bg-[#f59e0b] py-2 font-semibold text-[#030712] transition-colors hover:bg-[#d97706] disabled:cursor-not-allowed disabled:opacity-50"
					>{submitting ? 'Saving...' : 'Save Changes'}</button
				>
			</div>
		</form>
	</div>
</div>
