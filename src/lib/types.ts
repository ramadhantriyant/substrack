export interface Subscription {
	id: string;
	name: string;
	category: 'Entertainment' | 'SaaS' | 'Utilities' | 'Health' | 'Finance' | 'Other';
	amount: number;
	cycle: 'Monthly' | 'Yearly' | 'Weekly';
	nextBilling: string; // ISO date string YYYY-MM-DD
	icon: string;
	color: string;
	createdAt?: string;
}
