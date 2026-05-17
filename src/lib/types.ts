export interface Subscription {
	id: string;
	name: string;
	category: 'Entertainment' | 'SaaS' | 'Utilities' | 'Health' | 'Finance' | 'Other';
	amount: number;
	cycle: 'Monthly' | 'Yearly' | 'Quarterly';
	nextBilling: string; // ISO date string YYYY-MM-DD
	icon: string;
	color: string;
	currency?: string; // ISO 4217, e.g. 'USD'. Omitted/undefined → 'IDR'
	createdAt?: string;
}
