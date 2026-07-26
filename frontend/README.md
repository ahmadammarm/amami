# Amami - Frontend

The frontend of Amami is a modern, responsive, and highly interactive Single Page Application (SPA) designed for mosque administrators and community members.

## Technology Stack

- **Framework:** Vue 3 (Composition API, `<script setup>`)
- **Language:** TypeScript
- **Build Tool:** Vite
- **Styling:** Tailwind CSS 4
- **State Management:** Pinia
- **Routing:** Vue Router
- **Data Fetching:** Vue Query (`@tanstack/vue-query`), Axios
- **Form Handling:** VeeValidate + Zod
- **UI Components:** Reka UI (Headless UI), Unovis (Charts)
- **Icons:** Lucide Vue Next

## Architecture

The frontend follows an Atomic Design pattern for components:
- `src/components/atoms`: Basic elements (buttons, inputs)
- `src/components/molecules`: Small combinations of atoms (form fields with labels)
- `src/components/organisms`: Complex sections (sidebars, tables, modals)
- `src/components/pages`: Complete application views routed by Vue Router
- `src/components/templates`: Layout wrappers (DashboardLayout)

### API Layer

Network requests are abstracted into service files within `src/api/services/`. The application uses `@tanstack/vue-query` to manage server state, caching, and background synchronization, ensuring the UI remains perfectly in sync with the backend database.

### Strict Type Safety

The project utilizes TypeScript strictly. All API responses, request payloads, and Vue component props are defined in `src/types/` to prevent runtime errors and ensure code reliability. The `any` type is heavily discouraged.

## Getting Started

### Installation

Navigate to the `frontend` directory and install the dependencies using pnpm:

```bash
pnpm install
```

### Development Server

Start the Vite development server:

```bash
pnpm dev
```

### Building for Production

Compile TypeScript and build the optimized production bundle:

```bash
pnpm build
```

## RBAC & Router Guards

Navigation is strictly guarded based on the authenticated user's role. The `meta.requiresRole` property in `src/router/index.ts` dictates which roles are permitted to access specific pages. Unauthorized access attempts will redirect users to a Forbidden page or force a logout.
