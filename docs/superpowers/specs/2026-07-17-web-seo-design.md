# SEO для Nuxt web-приложения (twir.app)

Дата: 2026-07-17
Статус: утверждено пользователем (с правками: одна страница сравнения вместо `[slug]`)

## Цель

Улучшить индексацию и позиции twir.app в Google: выстроить технический SEO-фундамент и добавить
контентную страницу сравнения с конкурентами на всех 8 локалях.

## Контекст (текущее состояние)

- `web/` — Nuxt 4 (`compatibilityVersion: 4`), layers: `landing`, `dashboard`, `url-shortener`,
  `pastebin`, `public`, `overlays`.
- `@nuxtjs/seo` 3.4.0 уже установлен (включает robots, sitemap, og-image, schema-org, seo-utils).
- `@nuxtjs/i18n` 10.4.0, локали: `de, en, es, ja, pt, ru, sk, uk`, `defaultLocale: 'en'`.
  Стратегия по умолчанию (`prefix_except_default`).
- Уже есть: `site: { indexable: true }` (без `url`/`name`), robots с disallow
  (`/dashboard`, `/s`, `/h`, `/overlays`), базовые мета в `app/app.vue` со статичным
  og-image `/meta.webp`.
- Продакшен-домен: `twir.app` (подтверждено в `docker-compose.stack.yml`).
- Логотипы `nightbot.svg` и `streamelements.svg` уже существуют в
  `web/layers/dashboard/assets/integrations/` (неоптимизированы: 33КБ и 362КБ).

### Инвентарь роутов (актуально для sitemap/robots)

| Роут | Индексируем? |
| --- | --- |
| `/`, `/terms` (landing) | да |
| `/compare` (новый) | да |
| `/p/[channelName]/**` (публичные страницы каналов) | да |
| `/url-shortener` | да |
| `/login`, `/login/kick` | нет (utility) |
| `/url-shortener/profile`, `/h/profile` | нет (auth) |
| `/dashboard/**` | нет (robots + `ssr: false`) |
| `/s/**`, `/h/**` | нет (robots) |
| `/o/**`, `/overlays/**` | нет (robots, client-only) |

## Дизайн

### 1. Фундамент (`web/nuxt.config.ts`)

- `site`:
  - `url` — из env `NUXT_PUBLIC_SITE_URL`, дефолт `https://twir.app`
  - `name: 'Twir'`
  - `description` — как сейчас в `app.vue`
  - `defaultLocale: 'en'`
- `i18n.baseUrl` — та же env (включает корректные canonical/hreflang URL в `@nuxtjs/i18n`).
- `sitemap.exclude` — роуты из таблицы выше с пометкой «нет», включая локализованные варианты
  (проверить фактическую генерацию, см. «Верификация»).
- `schemaOrg.identity` — `defineOrganization`: name `Twir`, url, logo `/twir.svg`,
  `sameAs` (GitHub org, Discord invite, Twitter/X — взять актуальные ссылки из футера лендинга).
- `robots` — без изменений; модуль sitemap сам добавит директиву `Sitemap:`.

### 2. `app/app.vue` + командабл `useAppSeo`

- Новый `web/app/composables/use-app-seo.ts`: обёртка над `useSeoMeta` +
  `defineOgImageComponent('Twir', ...)`. API: `useAppSeo({ title, description })` —
  вызывается на страницах с локализованными строками.
- В `app.vue`: `useLocaleHead({ seo: true })` → hreflang (8 локалей + x-default), canonical,
  `og:locale` (+ alternate). Существующие `useHead`/`useSeoMeta` рефакторятся под `site`-конфиг
  (убрать дублирование description/keywords в один источник).

### 3. Главная страница (`layers/landing`)

- `pages/index.vue`: `useAppSeo` с локализованными title/description (новые ключи i18n).
- Якоря `id` у карточек фич в `components/index/features/` (например `id="commands"`,
  `id="song-requests"`) — для прямых ссылок и sitelinks.
- `defineSoftwareApplication` на главной: `applicationCategory: 'MultimediaApplication'`,
  `offers` (бесплатно, price 0), `featureList` из i18n.
- Новая FAQ-секция на главной (4–5 вопросов, i18n) + `defineQuestion`-schema (FAQPage) —
  кандидат на расширенный сниппет в Google.
- Упоминание страницы сравнения: компактный блок-ссылка в конце секции фич
  («Compare Twir with Nightbot, StreamElements…» → `/compare`) + пункт в футере
  (`layouts/default/footer.vue`).

### 4. Страница сравнения `/compare` (ОДНА страница)

- `layers/landing/pages/compare.vue` (без `[slug]`).
- Hero-блок: локализованный заголовок/подзаголовок («Twir vs другие Twitch-боты»).
- **Одна таблица**: колонки — Twir | Nightbot | StreamElements | Moobot | Fossabot;
  строки-фичи:
  - chat commands (custom commands)
  - timers
  - moderation
  - song requests
  - giveaways
  - overlays
  - chat alerts / notifications
  - games (!roulette и т.п.)
  - music integrations (Spotify, Last.fm)
  - open source / self-host
  - price
- Значения ячеек: `yes | partial | no` → иконки `Check` / `Minus` / `X` (lucide),
  `partial` — с пояснением (i18n note).
- Данные таблицы — типизированный data-файл
  `layers/landing/components/compare/compare.data.ts`: строки фич + значения по каждому боту,
  ключи i18n для названий/примечаний. Вёрстка таблицы — отдельный компонент
  `components/compare/compare-table.vue` (sticky header, адаптив: горизонтальный скролл на
  мобильных).
- **Логотипы в шапках колонок**: новая icon-коллекция `twir-compare`
  (`layers/landing/assets/compare/*.svg`, регистрация через `icon.customCollections` в
  `nuxt.config.ts`, как существующие `twir-integrations`/`twir-overlays`):
  - `nightbot.svg`, `streamelements.svg` — скопировать из
    `layers/dashboard/assets/integrations/`, оптимизировать через svgo
  - `moobot.svg`, `fossabot.svg` — взять из официальных источников (сайты/пресс-киты),
    привести к компактному SVG
  - логотипы без изменения фирменных цветов (nominative fair use)
- CTA-блок под таблицей: «Импорт из Nightbot/StreamElements» → `/dashboard/import`,
  «Попробовать Twir» → `/login`.
- SEO страницы: локализованные title/description, `defineOgImageComponent`, FAQ-блок
  (3–4 вопроса про переезд с конкурентов) + FAQPage schema, breadcrumb
  (`Home / Compare`) через `defineBreadcrumb`.
- i18n: namespace `compare` во всех 8 locale JSON. `en.json` — исходник, остальные 7
  переводятся в рамках задачи. `ja.json` сейчас почти пустой (36B) — ключи добавляем туда же,
  fallback на `en` уже обеспечен модулем.

### 5. OG-изображения

- Шаблон `web/app/components/OgImage/Twir.vue` (островной компонент nuxt-og-image): тёмный
  фон бренда (`#09090B`), логотип Twir, локализованные title/description, purple-gradient
  акцент как на лендинге.
- Главная, `/terms`, `/compare` используют шаблон через `useAppSeo`.

### 6. Верификация

1. `bun run build` в `web/` — сборка без ошибок.
2. Проверка артефактов: `/sitemap.xml` (все локали, отсутствие запрещённых роутов),
   `/robots.txt` (директива Sitemap).
3. Playwright: рендер `/compare` и `/ru/compare` — проверить `<head>`: hreflang (8 + x-default),
   canonical, og-теги (абсолютные URL), JSON-LD (Organization, FAQPage, Breadcrumb).
4. Lighthouse SEO-аудит главной и `/compare` — целевой score 100.
5. Ручная проверка og-image: открыть `/__og-image__/image/compare/og.png` (путь уточнить при
   реализации) — корректный рендер.

## Вне scope

- Отдельные страницы под каждую фичу (решено: якоря + schema на главной).
- Блог/гайды, публичная документация.
- Изменение robots-политики для `/p/**` и прочих уже индексируемых роутов.
- Контент самих дашборд-страниц (закрыты auth + `ssr: false`).

## Риски / заметки

- Товарные знаки: логотипы конкурентов используются для сравнения (nominative fair use), без
  искажений и без намёка на аффилиацию — в футере compare-страницы добавить дисклеймер.
- Точность данных таблицы — факты по конкурентам проверяются по их официальным сайтам на
  момент реализации; спорные ячейки помечаем `partial` с пояснением.
- og-image рендерится в рантайме (satori) — убедиться, что nitro preset `bun` тянет
  зависимости; иначе fallback — статичный `/meta.webp`.
