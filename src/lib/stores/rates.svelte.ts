const SYMBOLS = 'USD,EUR,GBP,SGD,JPY,AUD,CAD,MYR';
const CACHE_RATES = 'substrack_rates';
const CACHE_EXPIRY = 'substrack_rates_expiry';
const TTL_MS = 3 * 60 * 60 * 1000;

export const SUPPORTED_CURRENCIES = ['IDR', 'USD', 'EUR', 'GBP', 'SGD', 'JPY', 'AUD', 'CAD', 'MYR'] as const;

class RatesState {
	rates = $state<Record<string, number>>({});
	ready = $state(false);

	toIDR(amount: number, currency: string): number {
		if (currency === 'IDR') return amount;
		const rate = this.rates[currency];
		return rate ? amount / rate : amount;
	}

	async fetch(): Promise<void> {
		if (typeof localStorage === 'undefined') return;

		const expiry = localStorage.getItem(CACHE_EXPIRY);
		const cached = localStorage.getItem(CACHE_RATES);

		if (expiry && cached && Date.now() < Number(expiry)) {
			this.rates = JSON.parse(cached) as Record<string, number>;
			this.ready = true;
			return;
		}

		try {
			const res = await fetch(`https://api.frankfurter.dev/v1/latest?base=IDR&symbols=${SYMBOLS}`);
			if (!res.ok) throw new Error('fetch failed');
			const data = (await res.json()) as { rates: Record<string, number> };
			this.rates = data.rates;
			localStorage.setItem(CACHE_RATES, JSON.stringify(data.rates));
			localStorage.setItem(CACHE_EXPIRY, String(Date.now() + TTL_MS));
		} catch {
			if (cached) this.rates = JSON.parse(cached) as Record<string, number>;
		}

		this.ready = true;
	}
}

export const ratesState = new RatesState();
