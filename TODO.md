# Substrack TODO

## Authentication
- [ ] Add `JWTSecret` field to `AppConfig` and load it from env
- [ ] Implement `POST /auth/register` handler (validate email/password, hash password, create user)
- [ ] Implement `POST /auth/login` handler (verify password, return JWT)
- [ ] Add auth middleware to protect non-auth routes
- [ ] Implement refresh token flow (`MakeRefreshToken`/`HashRefreshToken` are already in `utils/auth.go`)
- [ ] Add `POST /auth/logout` endpoint

## User
- [ ] Add user model to `internal/models/`
- [ ] Add `GET /api/user/me` — return current authenticated user
- [ ] Add `PUT /api/user/me` — update name/email
- [ ] Add `PUT /api/user/me/password` — change password
- [ ] Add `DELETE /api/user/me` — delete account

## User-scoped resources
- [ ] Scope subscription endpoints to the authenticated user
- [ ] Scope category endpoints to the authenticated user
- [ ] Add `GET /api/user/me/subscription` — list user subscriptions
- [ ] Add `POST /api/user/me/subscription/{id}` — attach subscription to user
- [ ] Add `DELETE /api/user/me/subscription/{id}` — detach subscription from user
- [ ] Add `GET /api/user/me/category` — list user categories
- [ ] Add `POST /api/user/me/category/{id}` — attach category to user
- [ ] Add `DELETE /api/user/me/category/{id}` — detach category from user

## Done
- [x] Categories CRUD
- [x] Subscriptions CRUD
- [x] DB schema and migrations for users, users_categories, users_subscriptions
- [x] sqlc queries for users, users_categories, users_subscriptions
- [x] Auth utilities (JWT, argon2id password hashing, token extraction)
- [x] Auth routes registered (`/auth/login`, `/auth/register`)
- [x] Unique index on `users.email`
