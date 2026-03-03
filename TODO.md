# Substrack TODO

## User-scoped resources
- [ ] Scope subscription endpoints to the authenticated user
- [ ] Scope category endpoints to the authenticated user

## Done
- [x] Categories CRUD
- [x] Subscriptions CRUD
- [x] DB schema and migrations for users, users_categories, users_subscriptions
- [x] sqlc queries for users, users_categories, users_subscriptions
- [x] Auth utilities (JWT, argon2id password hashing, token extraction)
- [x] Unique index on `users.email`
- [x] `JWTSecret` in `AppConfig`, loaded from `JWT_SECRET` env var
- [x] `POST /auth/register` — validate, hash password, create user
- [x] `POST /auth/login` — verify password, return JWT + refresh token
- [x] `POST /auth/refresh` — issue new access token from refresh token
- [x] `POST /auth/logout` — revoke refresh token
- [x] Auth middleware (`RequireAuth`) protecting all non-public routes
- [x] User model in `internal/models/`
- [x] `GET /api/user/me`
- [x] `PUT /api/user/me`
- [x] `PUT /api/user/me/password`
- [x] `DELETE /api/user/me`
- [x] `GET /api/user/me/subscription`
- [x] `POST /api/user/me/subscription/{id}`
- [x] `DELETE /api/user/me/subscription/{id}`
- [x] `GET /api/user/me/category`
- [x] `POST /api/user/me/category/{id}`
- [x] `DELETE /api/user/me/category/{id}`
