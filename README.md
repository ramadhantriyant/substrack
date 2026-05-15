# substrack

A subscription tracker SPA — track every subscription, miss nothing.

Built with SvelteKit 5 (runes), Firebase, Tailwind CSS v4, and deployed on Cloudflare Workers.

## Stack

- **Frontend**: SvelteKit 5 (Svelte runes), Tailwind CSS v4, DM Sans / DM Mono
- **Auth**: Firebase Authentication (Google Sign-In)
- **Database**: Cloud Firestore (real-time listener)
- **Deployment**: Cloudflare Workers via `@sveltejs/adapter-cloudflare`
- **Runtime**: Bun

## Developing

Copy the environment file and fill in your Firebase credentials:

```sh
cp .env.example .env.local
```

Install dependencies and start the dev server:

```sh
bun install
bun run dev
```

## Building

```sh
bun run build
```

## Deploying

Authenticate with Cloudflare (one-time):

```sh
wrangler login
```

Build and deploy:

```sh
bun run deploy
```

After the first deploy, add your `*.workers.dev` domain to **Firebase Console → Authentication → Settings → Authorized domains** so Google Sign-In works in production.

## Environment Variables

| Variable | Description |
|---|---|
| `VITE_FIREBASE_API_KEY` | Firebase API key |
| `VITE_FIREBASE_AUTH_DOMAIN` | Firebase auth domain |
| `VITE_FIREBASE_PROJECT_ID` | Firestore project ID |
| `VITE_FIREBASE_STORAGE_BUCKET` | Firebase storage bucket |
| `VITE_FIREBASE_MESSAGING_SENDER_ID` | Firebase messaging sender ID |
| `VITE_FIREBASE_APP_ID` | Firebase app ID |

These are baked into the bundle at build time from `.env.local` (gitignored).
