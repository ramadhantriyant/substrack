<script lang="ts">
	import { subsState } from '$lib/stores/subscriptions.svelte';
	import { authState } from '$lib/stores/auth.svelte';
	import { signOut } from 'firebase/auth';
	import { auth } from '$lib/firebase';
	import { goto } from '$app/navigation';
	import { formatIDR, formatCurrency, toMonthly, getDaysUntil } from '$lib/utils';
	import { ratesState } from '$lib/stores/rates.svelte';
	import AddModal from '$lib/components/AddModal.svelte';
	import EditModal from '$lib/components/EditModal.svelte';
	import type { Subscription } from '$lib/types';
	import type { User } from 'firebase/auth';

	const CATEGORIES = ['All', 'Entertainment', 'SaaS', 'Utilities', 'Health', 'Finance', 'Other'];
	const CATEGORY_COLORS: Record<string, string> = {
		Entertainment: '#f59e0b',
		SaaS: '#6366f1',
		Utilities: '#10b981',
		Health: '#ec4899',
		Finance: '#3b82f6',
		Other: '#8b5cf6'
	};

	let showAddModal = $state(false);
	let editingSub = $state<Subscription | null>(null);
	let sortBy = $state<'urgency' | 'amount' | 'name'>('urgency');

	const uid = $derived(
		typeof authState.user === 'object' && authState.user ? (authState.user as User).uid : ''
	);

	const totalMonthlyIDR = $derived(
		subsState.subscriptions.reduce(
			(sum, s) => sum + toMonthly(s.amount, s.cycle, s.currency, ratesState.rates),
			0
		)
	);

	const sortedSubs = $derived.by(() => {
		const subs = [...subsState.filteredSubs];
		if (sortBy === 'urgency')
			return subs.sort((a, b) => getDaysUntil(a.nextBilling) - getDaysUntil(b.nextBilling));
		if (sortBy === 'amount')
			return subs.sort(
				(a, b) =>
					toMonthly(b.amount, b.cycle, b.currency, ratesState.rates) -
					toMonthly(a.amount, a.cycle, a.currency, ratesState.rates)
			);
		return subs.sort((a, b) => a.name.localeCompare(b.name));
	});

	const categoryBreakdown = $derived.by(() => {
		if (totalMonthlyIDR === 0) return [];
		const groups: Record<string, number> = {};
		for (const sub of subsState.subscriptions) {
			groups[sub.category] =
				(groups[sub.category] ?? 0) +
				toMonthly(sub.amount, sub.cycle, sub.currency, ratesState.rates);
		}
		return Object.entries(groups)
			.map(([cat, amount]) => ({ cat, amount, pct: (amount / totalMonthlyIDR) * 100 }))
			.sort((a, b) => b.amount - a.amount);
	});

	function urgencyBadge(days: number) {
		if (days <= 0) return { label: 'Overdue', color: '#ef4444', text: '#fca5a5' };
		if (days <= 3) return { label: `${days}d left`, color: '#f59e0b', text: '#fde68a' };
		if (days <= 10) return { label: `${days}d left`, color: '#3b82f6', text: '#93c5fd' };
		return { label: `${days}d left`, color: '#1e293b', text: '#64748b' };
	}

	async function handleDelete(subId: string) {
		if (!uid) return;
		await subsState.delete(uid, subId);
	}

	async function handleSignOut() {
		await signOut(auth);
		void goto('/login');
	}
</script>

<!-- Header -->
<header class="sticky top-0 z-10 border-b border-[#1e293b] bg-[#030712]/95 backdrop-blur-sm">
	<div class="mx-auto flex max-w-7xl items-center justify-between px-4 py-4 sm:px-6">
		<div class="flex items-center gap-3">
			<span class="text-xl font-bold text-[#f1f5f9] sm:text-2xl">
				<span class="text-[#f59e0b]">◎</span> substrack
			</span>
			{#if subsState.dueSoon.length > 0}
				<span
					class="flex items-center gap-1.5 rounded-full bg-[#f59e0b]/10 px-2.5 py-1 text-xs font-medium text-[#f59e0b]"
				>
					<span class="pulse-amber h-1.5 w-1.5 rounded-full bg-[#f59e0b]"></span>
					{subsState.dueSoon.length} due soon
				</span>
			{/if}
		</div>
		<div class="flex items-center gap-2">
			<button
				onclick={() => (showAddModal = true)}
				class="rounded-lg bg-[#f59e0b] px-4 py-2 text-sm font-semibold text-[#030712] transition-colors hover:bg-[#d97706]"
				>+ Add</button
			>
			<button
				onclick={handleSignOut}
				class="rounded-lg border border-[#1e293b] px-4 py-2 text-sm text-[#64748b] transition-colors hover:border-[#334155] hover:text-[#f1f5f9]"
				>Sign out</button
			>
		</div>
	</div>
</header>

<main class="mx-auto max-w-7xl px-4 py-6 sm:px-6">
	<!-- Summary cards -->
	<div class="mb-6 grid grid-cols-2 gap-3 lg:grid-cols-4 lg:gap-4">
		<div class="rounded-xl border border-[#1e293b] bg-[#0a0f1e] p-4 lg:p-5">
			<p class="text-xs text-[#64748b] sm:text-sm">Monthly Burn</p>
			<p class="mt-1 font-mono text-lg font-semibold text-[#f1f5f9] sm:text-2xl">
				{formatIDR(totalMonthlyIDR)}
			</p>
		</div>
		<div class="rounded-xl border border-[#1e293b] bg-[#0a0f1e] p-4 lg:p-5">
			<p class="text-xs text-[#64748b] sm:text-sm">Annual Spend</p>
			<p class="mt-1 font-mono text-lg font-semibold text-[#f1f5f9] sm:text-2xl">
				{formatIDR(totalMonthlyIDR * 12)}
			</p>
		</div>
		<div class="rounded-xl border border-[#1e293b] bg-[#0a0f1e] p-4 lg:p-5">
			<p class="text-xs text-[#64748b] sm:text-sm">Due This Week</p>
			<p class="mt-1 font-mono text-lg font-semibold text-[#f1f5f9] sm:text-2xl">
				{subsState.dueSoon.length}
			</p>
			{#if subsState.dueSoon.length > 0}
				<p class="mt-0.5 truncate text-xs text-[#f59e0b]">{subsState.dueSoon[0].name}</p>
			{/if}
		</div>
		<div class="rounded-xl border border-[#1e293b] bg-[#0a0f1e] p-4 lg:p-5">
			<p class="text-xs text-[#64748b] sm:text-sm">Active</p>
			<p class="mt-1 font-mono text-lg font-semibold text-[#f1f5f9] sm:text-2xl">
				{subsState.subscriptions.length}
			</p>
		</div>
	</div>

	<div class="flex flex-col gap-6 lg:flex-row">
		<!-- Main content -->
		<div class="min-w-0 flex-1">
			<!-- Category filter bar -->
			<div class="mb-4 flex flex-wrap gap-2">
				{#each CATEGORIES as cat (cat)}
					<button
						onclick={() => (subsState.activeCategory = cat)}
						class="rounded-full px-3 py-1 text-sm transition-colors {subsState.activeCategory ===
						cat
							? 'bg-[#f59e0b] font-medium text-[#030712]'
							: 'bg-[#0f172a] text-[#64748b] hover:text-[#f1f5f9]'}">{cat}</button
					>
				{/each}
			</div>

			<!-- Sort control -->
			<div class="mb-4 flex items-center gap-3">
				<span class="text-sm text-[#64748b]">Sort by</span>
				{#each [['urgency', 'Urgency'], ['amount', 'Amount'], ['name', 'Name']] as [val, label] (val)}
					<button
						onclick={() => (sortBy = val as typeof sortBy)}
						class="text-sm transition-colors {sortBy === val
							? 'font-medium text-[#f59e0b]'
							: 'text-[#64748b] hover:text-[#f1f5f9]'}">{label}</button
					>
				{/each}
			</div>

			<!-- Subscription list -->
			{#if subsState.subsLoading}
				<div class="flex items-center justify-center py-20">
					<div
						class="h-7 w-7 animate-spin rounded-full border-2 border-[#f59e0b] border-t-transparent"
					></div>
				</div>
			{:else if sortedSubs.length === 0}
				<div class="flex flex-col items-center justify-center py-20 text-center">
					<span class="mb-3 text-4xl">📭</span>
					<p class="text-[#64748b]">No subscriptions yet</p>
					<button
						onclick={() => (showAddModal = true)}
						class="mt-4 rounded-lg bg-[#f59e0b]/10 px-4 py-2 text-sm text-[#f59e0b] transition-colors hover:bg-[#f59e0b]/20"
						>Add your first subscription</button
					>
				</div>
			{:else}
				<div class="space-y-3">
					{#each sortedSubs as sub, i (sub.id)}
						{@const days = getDaysUntil(sub.nextBilling)}
						{@const badge = urgencyBadge(days)}
						<div
							class="fade-in-up group relative rounded-xl border border-[#1e293b] bg-[#0f172a] p-4 transition-transform hover:-translate-y-0.5"
							style="border-left: 3px solid {sub.color}; animation-delay: {i * 40}ms"
						>
							<div class="flex items-start gap-3">
								<span class="mt-0.5 text-2xl">{sub.icon}</span>
								<div class="min-w-0 flex-1">
									<div class="flex flex-wrap items-center gap-2">
										<span class="font-medium text-[#f1f5f9]">{sub.name}</span>
										<span
											class="rounded px-1.5 py-0.5 text-xs"
											style="background: {CATEGORY_COLORS[sub.category]}22; color: {CATEGORY_COLORS[
												sub.category
											]}">{sub.category}</span
										>
									</div>
									<div class="mt-1 flex flex-wrap items-center gap-x-2 text-sm text-[#64748b]">
										<span>{sub.cycle}</span>
										<span>·</span>
										<span>Next: {sub.nextBilling}</span>
									</div>
								</div>
								<div class="flex flex-col items-end gap-2 text-right">
									<div>
										<p class="font-mono font-medium text-[#f1f5f9]">
											{formatCurrency(sub.amount, sub.currency ?? 'IDR')}
										</p>
										{#if sub.cycle !== 'Monthly' || (sub.currency && sub.currency !== 'IDR')}
											<p class="font-mono text-xs text-[#64748b]">
												{formatIDR(
													toMonthly(sub.amount, sub.cycle, sub.currency, ratesState.rates)
												)}/mo
											</p>
										{/if}
									</div>
									<span
										class="rounded-full px-2 py-0.5 text-xs"
										style="background: {badge.color}22; color: {badge.text}">{badge.label}</span
									>
									<div class="flex gap-2 opacity-0 transition-all group-hover:opacity-100">
										<button
											onclick={() => (editingSub = sub)}
											aria-label="Edit {sub.name}"
											class="text-[#64748b] hover:text-[#f59e0b]"
										>
											<svg
												xmlns="http://www.w3.org/2000/svg"
												width="15"
												height="15"
												viewBox="0 0 24 24"
												fill="none"
												stroke="currentColor"
												stroke-width="1.75"
												stroke-linecap="round"
												stroke-linejoin="round"
											>
												<path d="M11 4H4a2 2 0 0 0-2 2v14a2 2 0 0 0 2 2h14a2 2 0 0 0 2-2v-7" />
												<path d="M18.5 2.5a2.121 2.121 0 0 1 3 3L12 15l-4 1 1-4 9.5-9.5z" />
											</svg>
										</button>
										<button
											onclick={() => handleDelete(sub.id)}
											aria-label="Delete {sub.name}"
											class="text-[#64748b] hover:text-red-400"
										>
											<svg
												xmlns="http://www.w3.org/2000/svg"
												width="15"
												height="15"
												viewBox="0 0 24 24"
												fill="none"
												stroke="currentColor"
												stroke-width="1.75"
												stroke-linecap="round"
												stroke-linejoin="round"
											>
												<path d="M3 6h18" />
												<path d="M19 6v14a2 2 0 0 1-2 2H7a2 2 0 0 1-2-2V6" />
												<path d="M8 6V4a2 2 0 0 1 2-2h4a2 2 0 0 1 2 2v2" />
												<line x1="10" y1="11" x2="10" y2="17" />
												<line x1="14" y1="11" x2="14" y2="17" />
											</svg>
										</button>
									</div>
								</div>
							</div>
						</div>
					{/each}
				</div>
			{/if}
		</div>

		<!-- Sidebars -->
		<div class="flex flex-col gap-4 lg:w-72 lg:shrink-0">
			<!-- Due Soon -->
			<div class="rounded-xl border border-[#1e293b] bg-[#0a0f1e] p-4">
				<h3 class="mb-3 text-sm font-semibold text-[#f1f5f9]">Due Soon</h3>
				{#if subsState.dueSoon.length === 0}
					<p class="text-sm text-[#64748b]">No upcoming bills in 7 days</p>
				{:else}
					<div class="space-y-3">
						{#each subsState.dueSoon as sub (sub.id)}
							{@const badge = urgencyBadge(getDaysUntil(sub.nextBilling))}
							<div class="flex items-center gap-2.5">
								<span class="text-xl">{sub.icon}</span>
								<div class="min-w-0 flex-1">
									<p class="truncate text-sm font-medium text-[#f1f5f9]">{sub.name}</p>
									<p class="font-mono text-xs text-[#64748b]">
										{formatCurrency(sub.amount, sub.currency ?? 'IDR')}
									</p>
								</div>
								<span
									class="shrink-0 rounded-full px-2 py-0.5 text-xs"
									style="background: {badge.color}22; color: {badge.text}">{badge.label}</span
								>
							</div>
						{/each}
					</div>
				{/if}
			</div>

			<!-- Category breakdown -->
			<div class="rounded-xl border border-[#1e293b] bg-[#0a0f1e] p-4">
				<h3 class="mb-3 text-sm font-semibold text-[#f1f5f9]">By Category</h3>
				{#if categoryBreakdown.length === 0}
					<p class="text-sm text-[#64748b]">No data yet</p>
				{:else}
					<div class="space-y-3">
						{#each categoryBreakdown as item (item.cat)}
							<div>
								<div class="mb-1 flex items-center justify-between text-xs">
									<span class="text-[#94a3b8]">{item.cat}</span>
									<span class="font-mono text-[#64748b]">{item.pct.toFixed(0)}%</span>
								</div>
								<div class="h-1.5 w-full overflow-hidden rounded-full bg-[#1e293b]">
									<div
										class="h-full rounded-full transition-all"
										style="width: {item.pct}%; background: {CATEGORY_COLORS[item.cat]}"
									></div>
								</div>
							</div>
						{/each}
					</div>
				{/if}
			</div>
		</div>
	</div>
</main>

{#if showAddModal}
	<AddModal onClose={() => (showAddModal = false)} {uid} />
{/if}

{#if editingSub}
	<EditModal onClose={() => (editingSub = null)} {uid} subscription={editingSub} />
{/if}
