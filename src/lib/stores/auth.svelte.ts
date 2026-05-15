import { onAuthStateChanged } from 'firebase/auth';
import { auth } from '$lib/firebase';
import type { User } from 'firebase/auth';

class AuthState {
	user = $state<User | null | false>(null);

	get isLoading() {
		return this.user === null;
	}

	get isAuthed() {
		return this.user !== null && this.user !== false;
	}
}

export const authState = new AuthState();

onAuthStateChanged(auth, (u) => {
	authState.user = u ?? false;
});
