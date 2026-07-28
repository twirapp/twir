# Twir Web Design System

## 1. Atmosphere & Identity

Twir is a compact operational dashboard. Neutral tonal surfaces, semantic controls, and platform icons make live status information scannable without introducing platform-specific layout patterns.

## 2. Color

| Role | Token | Light | Dark | Usage |
| --- | --- | --- | --- | --- |
| Background | `--background` | `oklch(1 0 0)` | `oklch(0.145 0 0)` | Page canvas |
| Surface | `--card` | `oklch(1 0 0)` | `oklch(0.205 0 0)` | Cards and panels |
| Text | `--foreground` | `oklch(0.145 0 0)` | `oklch(0.985 0 0)` | Primary content |
| Secondary text | `--muted-foreground` | `oklch(0.556 0 0)` | `oklch(0.708 0 0)` | Descriptions and metadata |
| Interactive | `--primary` | `oklch(0.205 0 0)` | `oklch(0.922 0 0)` | Primary actions and focus |
| Error | `--destructive` | `oklch(0.577 0.245 27.325)` | `oklch(0.704 0.191 22.216)` | Error feedback |

Use semantic Tailwind utilities backed by these tokens. Platform identity is communicated through its Simple Icons mark, not a custom color system.

## 3. Typography

The application uses `Inter, system-ui, Avenir, Helvetica, Arial, sans-serif`. Existing dashboard hierarchy uses `text-xl font-semibold` for section headings, `text-sm font-medium` for card titles, and `text-sm text-muted-foreground` for supporting copy.

## 4. Spacing & Layout

Spacing uses Tailwind's 4px scale. Dashboard action groups are vertical stacks with `gap-4`; cards use the shared `Card` padding and a responsive `flex` cluster for inline status and actions. Popup callbacks use a bounded `min-h-[100dvh]` cover with one centered card.

## 5. Components

### Action Card

- **Structure**: `Card` with `CardHeader` and `CardContent`.
- **States**: checking, disconnected, connected, opening authorization, replacement confirmation, callback error.
- **Accessibility**: buttons retain native semantics; replacement uses `ActionConfirm`; error feedback uses `Alert`.
- **Layout**: stack on the admin page; action/status cluster in card content.

### Bot Status Menu

- **Structure**: trigger button with platform icons and a dropdown of individual bindings.
- **States**: loading, online, disabled, mutation pending.
- **Accessibility**: keyboard-accessible dropdown menu with disabled pending actions.
- **Layout**: compact inline cluster in the persistent header.

## 6. Motion & Interaction

Shared components provide short transform/opacity transitions. Loading uses the existing spinner animation; no layout properties animate. Popups close after successful authorization.

## 7. Depth & Surface

The system uses the shared card treatment: semantic card surfaces, a 1px `--border`, rounded `--radius`, and `shadow-sm` for elevation. Do not add bespoke panel shadows or arbitrary platform-colored backgrounds.

## 8. Accessibility Constraints & Accepted Debt

Target WCAG 2.2 AA. All actions must have visible text labels, be keyboard reachable, preserve visible focus styles, and show authorization failures in persistent content as well as a toast. No accepted debt.
