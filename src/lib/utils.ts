export function formatIDR(amount: number): string {
	return new Intl.NumberFormat('id-ID', {
		style: 'currency',
		currency: 'IDR',
		maximumFractionDigits: 0
	}).format(amount);
}

export function toMonthly(amount: number, cycle: string): number {
	if (cycle === 'Weekly') return (amount * 52) / 12;
	if (cycle === 'Yearly') return amount / 12;
	return amount;
}

export function getDaysUntil(dateStr: string): number {
	const today = new Date();
	today.setHours(0, 0, 0, 0);
	const target = new Date(dateStr + 'T00:00:00');
	return Math.round((target.getTime() - today.getTime()) / (1000 * 60 * 60 * 24));
}
