import {
	collection,
	addDoc,
	deleteDoc,
	updateDoc,
	doc,
	onSnapshot,
	query,
	orderBy,
	Timestamp
} from 'firebase/firestore';
import { db } from '$lib/firebase';
import { toMonthly, getDaysUntil } from '$lib/utils';
import type { Subscription } from '$lib/types';
import type { User } from 'firebase/auth';

class SubscriptionState {
	subscriptions = $state<Subscription[]>([]);
	activeCategory = $state<string>('All');
	subsLoading = $state<boolean>(false);

	get filteredSubs() {
		return this.activeCategory === 'All'
			? this.subscriptions
			: this.subscriptions.filter((s) => s.category === this.activeCategory);
	}

	get totalMonthly() {
		return this.subscriptions.reduce((sum, s) => sum + toMonthly(s.amount, s.cycle), 0);
	}

	get dueSoon() {
		return this.subscriptions
			.filter((s) => getDaysUntil(s.nextBilling) <= 7)
			.sort((a, b) => getDaysUntil(a.nextBilling) - getDaysUntil(b.nextBilling));
	}

	#unsubscribe: (() => void) | null = null;

	startListener(user: User) {
		if (this.#unsubscribe) return;
		this.subsLoading = true;
		const q = query(
			collection(db, 'users', user.uid, 'subscriptions'),
			orderBy('nextBilling', 'asc')
		);
		this.#unsubscribe = onSnapshot(q, (snapshot) => {
			this.subscriptions = snapshot.docs.map((docSnap) => {
				const data = docSnap.data();
				const nextBilling =
					data.nextBilling instanceof Timestamp
						? data.nextBilling.toDate().toISOString().split('T')[0]
						: data.nextBilling;
				return { id: docSnap.id, ...data, nextBilling } as Subscription;
			});
			this.subsLoading = false;
		});
	}

	stopListener() {
		this.#unsubscribe?.();
		this.#unsubscribe = null;
		this.subscriptions = [];
	}

	async add(uid: string, sub: Omit<Subscription, 'id'>) {
		const { nextBilling, ...rest } = sub;
		await addDoc(collection(db, 'users', uid, 'subscriptions'), {
			...rest,
			nextBilling: Timestamp.fromDate(new Date(nextBilling + 'T00:00:00')),
			createdAt: new Date().toISOString()
		});
	}

	async delete(uid: string, subId: string) {
		await deleteDoc(doc(db, 'users', uid, 'subscriptions', subId));
	}

	async update(uid: string, subId: string, data: Partial<Omit<Subscription, 'id'>>) {
		await updateDoc(doc(db, 'users', uid, 'subscriptions', subId), data);
	}
}

export const subsState = new SubscriptionState();
