Build a subscription tracker web app called SubLens using SvelteKit and Firebase.

## Tech Stack
- SvelteKit with TypeScript and static adapter (SPA mode)
- Firebase Authentication (Google Sign-In)
- Cloud Firestore (real-time listener)
- Tailwind CSS for styling

## Firebase Setup
Create src/lib/firebase.ts that initializes Firebase using these environment variables:
- VITE_FIREBASE_API_KEY
- VITE_FIREBASE_AUTH_DOMAIN
- VITE_FIREBASE_PROJECT_ID
- VITE_FIREBASE_STORAGE_BUCKET
- VITE_FIREBASE_MESSAGING_SENDER_ID
- VITE_FIREBASE_APP_ID

Export `auth` and `db` from this file.

## TypeScript Types — src/lib/types.ts
```ts
export interface Subscription {
  id: string;
  name: string;
  category: "Entertainment" | "SaaS" | "Utilities" | "Health" | "Finance" | "Other";
  amount: number;
  cycle: "Monthly" | "Yearly" | "Weekly";
  nextBilling: string; // ISO date string YYYY-MM-DD
  icon: string;
  color: string;
  createdAt?: string;
}
```

## Svelte Stores — src/lib/stores/

### auth.ts
- `user` writable store with three states: null (loading), false (not authed), User (authed)
- Initialize onAuthStateChanged listener at module level
- Export derived stores: `isLoading`, `isAuthed`

### subscriptions.ts
- `subscriptions` writable store — full list from Firestore
- `activeCategory` writable store — default "All"
- `filteredSubs` derived store — filters subscriptions by activeCategory
- `totalMonthly` derived store — sum of all subscriptions converted to monthly, derived from full subscriptions (not filteredSubs)
- `dueSoon` derived store — subscriptions where nextBilling is within 7 days, sorted ascending
- `subsLoading` writable boolean store
- `startSubscriptionListener(user)` — attaches onSnapshot listener ordered by nextBilling asc, converts Firestore Timestamps to ISO date strings
- `stopSubscriptionListener()` — detaches listener and clears store
- `addSubscription(uid, sub)` — addDoc to users/{uid}/subscriptions, store nextBilling as Firestore Timestamp
- `deleteSubscription(uid, subId)` — deleteDoc
- `updateSubscription(uid, subId, data)` — updateDoc

## Routes

### src/routes/+layout.svelte
- Import user and isLoading from auth store
- Import startSubscriptionListener, stopSubscriptionListener from subscriptions store
- Reactively watch user store: if user is a User object call startSubscriptionListener; if false call stopSubscriptionListener and redirect to /login; if null show a full-screen loading state
- Use $page.url.pathname to avoid redirect loop on /login

### src/routes/login/+page.svelte
- Full-page centered login screen
- Dark theme matching the dashboard
- App logo and tagline: "Track every subscription. Miss nothing."
- Single "Continue with Google" button using signInWithPopup with GoogleAuthProvider
- On success redirect to /dashboard
- If user is already authed, redirect to /dashboard on mount

### src/routes/dashboard/+page.svelte
- Redirect to /login if not authed
- Import filteredSubs, activeCategory, totalMonthly, dueSoon, subsLoading from subscriptions store
- Import user from auth store

Header:
- App logo (◎ SubLens)
- Pulsing amber dot badge showing count if dueSoon has items
- "Add" button that opens AddModal
- Sign out button

Summary cards (4 cards in a responsive grid):
- Monthly Burn: totalMonthly formatted as IDR
- Annual Spend: totalMonthly * 12 formatted as IDR
- Due This Week: count of dueSoon, subtitle shows next due subscription name
- Active Subscriptions: total count

Category filter bar:
- Categories: ["All", "Entertainment", "SaaS", "Utilities", "Health", "Finance", "Other"]
- Clicking a category calls activeCategory.set(cat)
- Active category highlighted in amber

Sort control:
- Options: Urgency (by nextBilling asc), Amount (by monthly equivalent desc), Name (alphabetical)
- Sorting is done on filteredSubs in the component

Subscription list:
- Renders filteredSubs using {#each $filteredSubs as sub (sub.id)}
- Each card shows: icon, name, category badge, billing cycle, next billing date, amount, monthly equivalent, urgency countdown badge
- Urgency badge: red for today/overdue, amber for ≤3 days, blue for ≤10 days, gray otherwise
- Delete button on each card

Due Soon sidebar panel:
- Shows subscriptions due within 7 days
- Each row: icon, name, amount, urgency badge

Category breakdown sidebar panel:
- Bar chart per category showing percentage of totalMonthly
- Color coded per category

## AddModal component — src/lib/components/AddModal.svelte
Fields: name (text), amount (number, IDR), cycle (select), category (select), nextBilling (date), icon (picker from preset list), color (color picker from preset swatches)
- On submit call addSubscription from store
- Validate name and amount are not empty
- Close on backdrop click or × button

## Utility functions — src/lib/utils.ts
- `formatIDR(amount)` — Intl.NumberFormat id-ID currency IDR
- `toMonthly(amount, cycle)` — converts to monthly equivalent
- `getDaysUntil(dateStr)` — days from today to date string

## Design
Dark financial-grade theme:
- Background: #030712
- Card background: #0a0f1e / #0f172a
- Accent: amber #f59e0b
- Borders: #1e293b
- Text primary: #f1f5f9
- Text muted: #64748b
- Fonts: DM Sans (body), DM Mono (numbers) from Google Fonts
- Category colors: Entertainment #f59e0b, SaaS #6366f1, Utilities #10b981, Health #ec4899, Finance #3b82f6, Other #8b5cf6
- Smooth hover transitions on cards (translateY -2px)
- Left border on cards colored by subscription color
- Staggered fade-in animation on list items

## svelte.config.js
Use @sveltejs/adapter-static with fallback: "index.html" for SPA mode.

## .env.example
Include all VITE_FIREBASE_* keys as empty placeholders.

## Notes
- After addSubscription or deleteSubscription, do NOT manually update the store — the Firestore onSnapshot listener handles it automatically
- totalMonthly must derive from the full subscriptions store, not filteredSubs, so summary cards always show total spend regardless of active category filter
- The subscriptions listener should be attached once on login and detached on logout