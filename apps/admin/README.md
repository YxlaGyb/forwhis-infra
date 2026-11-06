# Admin

Nuxt 4 + Vue 3 + Ant Design Vue administration app. Built as a static SPA.

## Setup

```bash
pnpm install
```

Copy `.env.example` to `.env` and point `NUXT_PUBLIC_API_BASE_URL` at the Go API.

## Development

```bash
pnpm run dev
```

Runs at <http://localhost:3000>. From the repository root, `pnpm dev:admin` does the same.

## Build

```bash
pnpm run build
```

Produces a static site in `.output/public`. Preview it with `pnpm run preview`.

`NUXT_PUBLIC_*` values are inlined at build time. Build with the production API
URL (`NUXT_PUBLIC_API_BASE_URL=https://api.banana404.top`) or without a local
`.env`, otherwise the local development URL ends up in the bundle.

The host must serve `200.html` for unknown paths so client-side routes such as
`/editor/12` and `/settings` resolve. For example, Cloudflare Pages needs
`/* /index.html 200`.

## Checks

```bash
pnpm run typecheck
pnpm run lint
pnpm run format
```
