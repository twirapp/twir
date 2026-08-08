# План: прокачка конструктора оверлеев (по мотивам hepega studio)

Дата: 2026-08-06
Референс: https://hepega.com/studio (изучен через Playwright, см. раздел «Что у hepega»)

## Цель

Довести наш конструктор (`overlay-builder`) до уровня hepega: новые типы слоёв
(видео, iframe/виджет, YouTube, текст, эмоции, RTE), удобная логика скрытия,
готовые виджеты Twir как слои. **Без источника «Экран»** (screen capture не делаем).

## Что у hepega (зафиксировано)

- **Типы слоёв**: Картинка, Видео, Эмоции, YouTube, Текст, Виджет (=iframe с URL,
  масштабом и кнопкой «выключить в редакторе для экономии ОЗУ»), Экран, RTE
  (саб-меню: Чат, Файтинг эмоутов — т.е. готовые встроенные виджеты).
- **hewidget** (их игровые виджеты как слои): heguess, hepoll, hebingo, hecodil,
  heless, heafk, hepredict.
- **Эмоции**: модалка с табами `7TV` (глобальные), `7TV канал`, `Twitch`
  (глобальные), `Сабки канала`, поиск по имени/тегу, грид эмоутов.
- **На слое** (флоатинг-тулбар): прозрачность, масштаб, скрыть/показать,
  отразить по горизонтали, анимация по траектории, lock, удалить, поворот
  (Shift = снап 15°).
- **Логика скрытия**: у каждого элемента бейдж «не на стриме», кнопка «Скрыть»
  в тулбаре и в списке элементов, глобальный тумблер «На показ» = «добавлять
  контент выключенным». Скрытие realtime-синкается в OBS.
- **Прочее**: магнит к сетке (Alt = временно off), линейки-направляющие, сцены,
  история undo/redo, хоткеи (стрелки 5px/1px, Shift+drag marquee, Alt+resize =
  crop как в OBS, Shift+resize = aspect ratio), звуковые слои, выбор разрешения
  зоны стрима, OBS-ссылка.

## Что у нас уже есть (по explore)

- Редактор: `web/layers/dashboard/features/overlay-builder/` (OverlayBuilder.vue,
  `components/Canvas.vue` на `vue3-moveable` — drag/resize/rotate/снаппинг/гайды
  уже есть, `components/LayersPanel.vue`, `components/PropertiesPanel.vue`).
- Слой: `Layer { type, posX, posY, width, height, rotation, opacity, visible,
  locked, zIndex, settings }` — **visible/locked/opacity уже персистятся**,
  рантайм фильтрует `layer.visible`.
- Типы сейчас: только `HTML` и `IMAGE` (GraphQL enum `ChannelOverlayLayerType`,
  entity `libs/entities/custom_overlay`, model
  `libs/repositories/channels_overlays/model`, PG enum
  `channels_overlays_layers_type`, settings — JSONB).
- CRUD: GraphQL `channelOverlayCreate/Update` (whole-overlay replace, max 15
  слоёв), подписка `customOverlaySettings` для рантайма, instant-save позиций.
- Рантайм: `frontend/overlays/src/pages/overlays.vue` — dispatch
  `html-layer.vue` / `image-layer.vue`.
- Эмоуты: 7TV/BTTV/FFZ инфраструктура есть (`libs/integrations/seventv`,
  `apps/emotes-cacher`, рантайм-composables `use-seven-tv` и т.д.), но **готового
  emote-picker компонента в дашборде нет**.
- Виджеты: predictions (eventsub + parser + миграция), polls (eventsub
  handler), chat-оверлей (отдельная страница), giveaways. Guess/bingo/afk —
  отдельных модулей нет.

## Дизайн-решения

1. **Новые типы слоёв** — расширяем существующий путь
   `enum → entity/model → миграция → builder menu → properties panel → canvas
   preview → runtime component`. Не плодим отдельные таблицы: всё в `settings`
   JSONB, дискриминация по `type`.
2. **RTE = «Виджеты Twir»**: слой-iframe, который вместо ручного URL даёт выбор
   из реестра наших оверлеев (chat, и далее predictions/polls по мере появления
   их overlay-страниц) и сам подставляет публичный URL + apiKey. Это и есть
   «связать с твировскими» без дублирования рендера внутри конструктора.
   heguess/hebingo/hecodil/heless/heafk — отдельные игровые модули, в этот план
   не входят; когда появятся, они встают в тот же реестр.
3. **Эмоции** — отдельный тип `EMOTE` (один эмоут на слой, как у hepega) +
   переиспользуемый компонент `EmotePicker` (пригодится и в других местах
   дашборда).
4. **Логика скрытия** — докручиваем UX поверх существующего `visible`:
   тумблер «добавлять скрытыми», бейдж «скрыт на стриме» в списке слоёв,
   кнопка show/hide в тулбаре слоя; instant-save уже realtime-развозит по
   подписке. Триггеров/расписаний пока нет.

## Этапы

### Этап 1 — базовые типы: TEXT, VIDEO, IFRAME, YOUTUBE

Backend:
- Миграция: `bun cli m create --name add_overlay_layer_types --db postgres --type sql`
  → `ALTER TYPE channels_overlays_layers_type ADD VALUE ...` для TEXT, VIDEO,
  IFRAME, YOUTUBE.
- `libs/entities/custom_overlay/custom_overlay.go` — новые константы типов.
- `libs/repositories/channels_overlays/model/model.go` — расширить
  `OverlayLayerSettings`: `textContent`, `textFont`, `textSize`, `textColor`,
  `textAlign`, `videoUrl`, `videoLoop`, `videoMuted`, `iframeUrl`, `iframeScale`,
  `iframeDisabledInEditor`, `youtubeVideoId`, `youtubeAutoplay`, `youtubeLoop`.
- GraphQL: `apps/api-gql/internal/delivery/gql/schema/overlays/overlays-custom.graphql`
  — enum + поля в `ChannelOverlayLayerSettings`/`Input`; маппер
  `internal/delivery/gql/mappers/channel_overlay.go`; `bun cli build gql`.

Frontend builder:
- `features/overlay-builder/types/index.ts` — типы и settings-поля.
- `features/overlay-builder/OverlayBuilder.vue` — меню добавления: группы
  «Медиа» (Картинка, Видео, YouTube), «Контент» (Текст, HTML, Виджет/iframe,
  Эмоции), «Виджеты Twir».
- `features/overlay-builder/components/PropertiesPanel.vue` — редакторы свойств
  под каждый тип.
- `features/overlay-builder/components/Canvas.vue` — preview-компоненты слоёв.

Runtime `frontend/overlays`:
- `components/text-layer.vue` (styled div), `video-layer.vue` (`<video>`,
  muted+autoplay для OBS), `iframe-layer.vue` (iframe + scale),
  `youtube-layer.vue` (youtube-nocookie embed, autoplay/loop/playlist=params).
- `pages/overlays.vue` — dispatch по новым типам; расширить `Layer` interface в
  `composables/overlays/use-overlays.ts` (поля settings уже прилетают).

### Этап 2 — Эмоции (EMOTE + EmotePicker)

- Тип `EMOTE`, settings: `emoteUrl`, `emoteName`, `emoteProvider`
  (7TV/TWITCH), `flipH`.
- Компонент `web/layers/dashboard/components/emote-picker/EmotePicker.vue`:
  - Табы: 7TV (глобальный сет — 7TV REST v3 `emote-sets/global`), 7TV канала
    (`emote-sets` по twitch id канала), Twitch глобальные, сабы канала.
  - Twitch-эмоуты: Helix `chat/emotes` требует токен — делаем через наш
    api-gql (новый query, токен канала уже есть в tokens-сервисе) либо через
    emotes-cacher данные; решить на реализации, предпочтительно api-gql.
  - Поиск по имени, грид с ленивой подгрузкой.
- Рантайм `emote-layer.vue` — по сути image-layer + flip; можно переиспользовать.

### Этап 3 — «Виджеты Twir» (аналог RTE) и UX скрытия

- Реестр виджетов в дашборде: `{ key, name, urlBuilder(channelId, apiKey),
  settingsComponent }`. Стартовый набор: Chat overlay. Слой `IFRAME` с выбором
  «из реестра» или «свой URL».
- **Инлайн-настройки виджета без редиректа**: у записи реестра может быть
  `settingsComponent` — существующая форма настроек этого виджета
  (например `web/layers/dashboard/pages/dashboard/overlays/chat/components/Form.vue`
  для чата). В конструкторе по кнопке «Настроить» на слое открываем
  `DialogOrSheet`-модалку и рендерим этот компонент внутри — дашборд и
  конструктор живут в одном Nuxt-приложении, composables и мутации формы
  работают как есть. Форма сохраняет настройки виджета через его обычные
  GraphQL-мутации, а слой хранит только ссылку на виджет (`widgetKey`), так что
  изменения подхватываются iframe'ом автоматически. Требование к формам:
  не зависеть от route params страницы — если зависит (id из URL), оборачиваем
  в адаптер, который прокидывает нужные props.
- UX скрытия: тумблер «добавлять скрытыми» в тулбаре конструктора, бейдж
  «скрыт на стриме» в LayersPanel, кнопка show/hide на выделенном слое.
- По желанию (если быстро): flip-horizontal у слоёв (поле `flipH` в
  `channels_overlays_layers`, миграция) — у hepega есть, копеечная фича.

### Не входит

- Источник «Экран» (screen capture).
- Игровые виджеты heguess/hepoll/hebingo/hecodil/heless/heafk/hepredict как
  слои — ждут своих модулей; predictions/polls подключим через реестр, когда у
  них появятся overlay-страницы.
- Анимации по траектории, звуковые слои, сцены, линейки — отдельными задачами.

## Проверка

Общие команды после каждого этапа:
- `bun cli build gql` — регенерация проходит без ошибок (после правок схемы).
- `go build ./...` в `apps/api-gql`, `libs/entities`, `libs/repositories` —
  компиляция чистая.
- `bun lint` — без новых ошибок.
- `bun run build` в `frontend/overlays` — сборка exit 0.
- Typecheck dashboard (команда из `web/package.json`, напр. `bun run typecheck`
  или `nuxi typecheck`) — без ошибок в изменённых файлах.

QA по этапам (ручной прогон через дашборд + рантайм-ссылку оверлея):

Этап 1:
- TEXT: добавить слой → задать текст/размер/цвет/выравнивание → сохранить →
  в рантайме виден текст с заданными стилями; перезагрузка редактора —
  настройки на месте.
- VIDEO: слой с mp4-URL → в рантайме видео играет (muted autoplay), loop
  работает; в редакторе preview не душит вкладку.
- IFRAME: слой с произвольным URL → в рантайме iframe грузится; масштаб
  применяется; «выключить в редакторе» скрывает iframe только в редакторе.
- YOUTUBE: слой с video id → в рантайме играет embed, loop/autoplay по
  настройкам.
- Скрытие: toggle visible → в рантайме слой исчезает/появляется без
  перезагрузки (подписка), z-order сохраняется.

Этап 2:
- EMOTE: открыть EmotePicker → табы 7TV/7TV канал/Twitch/Сабки грузят сетки →
  поиск фильтрует → выбор эмоута ставит его на слой → в рантайме эмоут
  отображается; пересохранение/перезагрузка не ломает слой.

Этап 3:
- Виджет из реестра (Chat): добавить слой → URL чата подставился → «Настроить»
  открывает модалку с формой чата → изменение настроек сохраняется обычной
  мутацией чата → в рантайме чат работает внутри конструкторного оверлея.
