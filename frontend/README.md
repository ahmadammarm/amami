# Amami — Frontend

Single-page app for amami: the public mosque page, sign-in, and the committee dashboard. Vue 3 with TypeScript, built by Vite. All interface copy is in Bahasa Indonesia.

## Stack

| Concern | Library |
|---|---|
| Framework | Vue 3 (Composition API, `<script setup>`), TypeScript (strict) |
| Build | Vite 8 |
| Styling | Tailwind CSS 4 with project design tokens |
| UI primitives | shadcn-vue on Reka UI, icons from Lucide |
| Server state | TanStack Vue Query over an Axios client |
| Client state | Pinia (auth session) |
| Routing | Vue Router |
| Forms | VeeValidate + Zod |
| Charts | ApexCharts (`vue3-apexcharts`) |

## Getting started

**Requirements:** Node.js 24 (pinned in `.nvmrc`; `nvm use` or `fnm use` picks it up) and pnpm. The API from `../backend` should be running.

```bash
pnpm install
pnpm dev       # http://localhost:5173
pnpm build     # type-check (vue-tsc) and production build into dist/
pnpm preview   # serve the production build locally
```

The API base URL defaults to `http://localhost:8080/api/v1`. To point elsewhere, create `.env.local`:

```
VITE_API_BASE_URL=https://api.example.org/api/v1
```

`pnpm build` is the type check; there is no separate lint or test script.

## Project layout

```
src/
  api/          client.ts (request<T>, errors, asset URLs), queryKeys.ts, services/<module>.ts
  components/
    atoms/      shadcn-vue primitives and the mihrab-arch atoms (ArchBadge, ArchFrame)
    molecules/  PageHeader, StatCard, DataPanel, PaginationBar, EmptyState, StatusBadge,
                FormDialog, FormField, and module dialogs
    organisms/  AppSidebar, FeedbackDialog, ErrorBoundary
    pages/      One component per route
    templates/  DashboardLayout, AuthLayout, StatusLayout
  composables/  useFeedback and shared queries
  lib/          permissions, formatters, role labels, mustahik scoring helpers, service status
  schemas/      Zod schemas, one file per module; payload types are z.infer of these
  stores/       auth session
  types/        API response types
  style.css     design tokens and shared component classes
```

## How things fit together

- **API access.** Services call `request<T>()`, which unwraps the backend's `{ success, message, data, errors }` envelope. Query keys always come from `src/api/queryKeys.ts`; after a mutation, invalidate the module's `all` key.
- **Forms.** Each form binds a Zod schema from `src/schemas/` with VeeValidate (`useForm` + `toTypedSchema`), renders inside `FormDialog`, and shows field errors through `FormField`.
- **Notifications.** There are no toasts. `feedback.success`, `feedback.error` and `feedback.confirm` from `src/composables/useFeedback.ts` show queued dialogs; `confirm` resolves to `true` or `false` and replaces `window.confirm`.
- **Access control.** `/auth/me` returns the user's permission codes. Routes declare `meta.permission`, the sidebar uses the same codes, and pages show write actions only when `authStore.can(...)` allows them. Role names are for display only.
- **Status pages.** Unknown URLs render the 404 page. When the API returns 503, cannot be reached, or the device is offline, the app opens `/pemeliharaan`, retries every 15 seconds, and returns to the original page once the server answers. A stale session on a public page is dropped quietly instead of forcing a sign-in.

## Design system

Colours are named after the materials of a mosque and defined in `src/style.css`:

| Token | Use |
|---|---|
| `kapur` | App background (lime-plaster white) |
| `batu` | Borders, table header, input outlines |
| `jati` | Text (dark teak) |
| `sajadah` | Primary actions, sidebar (prayer-carpet green) |
| `kuningan` | Accents: active menu marker, focus ring (brass) |
| `bata` | Destructive actions and errors (brick) |

Type: Reem Kufi for headings (`font-display`), Plus Jakarta Sans for text, IBM Plex Mono for IDs and codes. Use these tokens instead of Tailwind's default palette, `.field-control` for inputs, `.tabel` for tables, and `.tombol-ikon` for icon-only buttons.
