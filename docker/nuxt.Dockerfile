ARG NODE_VERSION=26.7.0
ARG NPM_REGISTRY=https://registry.npmmirror.com

FROM docker.1ms.run/library/node:${NODE_VERSION}-alpine AS base

ARG NPM_REGISTRY
ENV npm_config_registry=${NPM_REGISTRY}
RUN npm install -g pnpm@11.22.0

WORKDIR /repo

FROM base AS deps

ARG NPM_REGISTRY

COPY package.json pnpm-lock.yaml pnpm-workspace.yaml ./
COPY apps/api/package.json apps/api/package.json
COPY apps/admin/package.json apps/admin/package.json

RUN --mount=type=cache,id=pnpm-store,target=/pnpm-store \
    --mount=type=cache,id=pnpm-metadata,target=/root/.cache/pnpm \
    pnpm install --frozen-lockfile --ignore-scripts --prefer-offline --store-dir=/pnpm-store \
    --network-concurrency=4 --fetch-retries=5 --fetch-timeout=180000 --fetch-retry-maxtimeout=180000 \
    --filter "./apps/admin..." --registry=${NPM_REGISTRY}

FROM deps AS build

ARG NUXT_PUBLIC_API_BASE_URL=/
ENV NUXT_PUBLIC_API_BASE_URL=${NUXT_PUBLIC_API_BASE_URL}
ENV NODE_OPTIONS=--max-old-space-size=2048

COPY apps/admin apps/admin
RUN pnpm --filter "./apps/admin" exec nuxt build

FROM docker.1ms.run/library/node:${NODE_VERSION}-alpine AS run

ENV NODE_ENV=production HOST=0.0.0.0 NITRO_PORT=23325 NODE_OPTIONS=--max-old-space-size=384
WORKDIR /app

COPY --from=build /repo/apps/admin/.output ./.output

USER node
EXPOSE 23325
CMD ["node", ".output/server/index.mjs"]
