export function formatIDR(amount: number): string {
	return new Intl.NumberFormat('id-ID', {
		style: 'currency',
		currency: 'IDR',
		maximumFractionDigits: 0
	}).format(amount);
}

const CURRENCY_LOCALES: Record<string, string> = {
	IDR: 'id-ID',
	USD: 'en-US',
	EUR: 'de-DE',
	GBP: 'en-GB',
	SGD: 'en-SG',
	JPY: 'ja-JP',
	AUD: 'en-AU',
	CAD: 'en-CA',
	MYR: 'ms-MY'
};

export function formatCurrency(amount: number, currency: string): string {
	const locale = CURRENCY_LOCALES[currency] ?? 'en-US';
	try {
		return new Intl.NumberFormat(locale, {
			style: 'currency',
			currency,
			maximumFractionDigits: ['IDR', 'JPY'].includes(currency) ? 0 : 2
		}).format(amount);
	} catch {
		return `${currency} ${amount.toFixed(2)}`;
	}
}

export function toMonthly(
	amount: number,
	cycle: string,
	currency = 'IDR',
	rates: Record<string, number> = {}
): number {
	const amountIDR = currency === 'IDR' || !rates[currency] ? amount : amount / rates[currency];
	if (cycle === 'Weekly') return (amountIDR * 52) / 12;
	if (cycle === 'Yearly') return amountIDR / 12;
	return amountIDR;
}

export function getDaysUntil(dateStr: string): number {
	const today = new Date();
	today.setHours(0, 0, 0, 0);
	const target = new Date(dateStr + 'T00:00:00');
	return Math.round((target.getTime() - today.getTime()) / (1000 * 60 * 60 * 24));
}
