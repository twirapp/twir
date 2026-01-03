# Исправление проблемы с IMAGE слоями

## Проблема
Картинка сохраняется, но не отображается, показывает "failed to load image" на канвасе.

## Что было исправлено

### 1. Валидация пустых URL
Теперь компоненты не пытаются загрузить пустую строку как изображение.

### 2. Исправлена передача событий в PropertiesPanel
Была ошибка в коде, которая не позволяла обновлять `imageUrl`.

### 3. Добавлено подробное логирование
Для отладки добавлены `console.log` на всех этапах загрузки и сохранения данных.

## Как проверить, в чём проблема

### Шаг 1: Откройте консоль браузера (F12)

Перейдите на вкладку Console.

### Шаг 2: Создайте или откройте IMAGE слой

Вы увидите логи в консоли:
```
[edit.vue] Loading overlay data from API: {...}
[edit.vue] Converting layer from API: { type: 'IMAGE', imageUrl: '...', ... }
[OverlayBuilder] Loading layer 0: { type: 'IMAGE', imageUrl: '...', ... }
[Canvas] IMAGE layer mounted: layer-xxx imageUrl: '...'
[ImageLayerPreview] Component mounted with imageUrl: '...' hasValidUrl: true/false
```

### Шаг 3: Проверьте значение imageUrl

#### Вариант А: imageUrl = undefined или null
```
[ImageLayerPreview] imageUrl: undefined hasValidUrl: false
```
**Проблема**: GraphQL не возвращает поле `imageUrl`
**Решение**: Проверьте бэкенд, что поле `imageUrl` включено в резолвер

#### Вариант Б: imageUrl = пустая строка ''
```
[ImageLayerPreview] imageUrl: '' hasValidUrl: false
```
**Проблема**: Слой создан без URL
**Решение**: В панели свойств введите URL изображения или нажмите "Use Placeholder Image"

#### Вариант В: imageUrl есть, но изображение не загружается
```
[ImageLayerPreview] imageUrl: 'https://example.com/image.png' hasValidUrl: true
[ImageLayerPreview] Failed to load image: 'https://example.com/image.png'
```
**Возможные причины**:
- Неправильный URL
- CORS блокирует загрузку (сервер не разрешает cross-origin запросы)
- Изображение не существует по этому URL
- Проблема с сетью

**Решения**:
1. Проверьте URL - откройте его в новой вкладке браузера
2. Используйте изображения с серверов, которые разрешают CORS:
   - `https://via.placeholder.com/300x200`
   - `https://picsum.photos/300/200`
3. Проверьте вкладку Network в DevTools на наличие ошибок CORS

## Частые проблемы и решения

### "No image URL" в редакторе

**Решение**:
1. Выберите слой в списке
2. В панели Properties справа найдите "Image URL"
3. Введите URL изображения или нажмите "Use Placeholder Image"

### "Failed to load image"

**Причины**:
- Неверный формат URL
- CORS блокировка
- Изображение не существует
- Проблема с сетью

**Решение**:
1. Откройте URL в новой вкладке - если не открывается, URL неверный
2. Используйте CORS-friendly хостинги для изображений
3. Проверьте DevTools → Network на наличие ошибок

### Изображение показывается в редакторе, но не в оверлее

**Решение**:
1. Откройте консоль на странице оверлея (не дашборда)
2. Проверьте, есть ли ошибки загрузки
3. Проверьте, что GraphQL запрос в `use-custom-overlay.ts` включает `imageUrl`

### imageUrl не сохраняется в базу данных

**Проверьте**:
1. В консоли браузера найдите лог сохранения:
   ```
   [edit.vue] Layers to be saved: [{ settings: { imageUrl: '...' } }]
   ```
2. Если `imageUrl` есть в логе, проблема на бэкенде
3. Проверьте GraphQL резолвер и миграции базы данных

## Тестирование

После внесения изменений проверьте:

- [ ] Создать новый IMAGE слой → видно placeholder изображение
- [ ] Ввести свой URL → изображение загружается и показывается
- [ ] Ввести неверный URL → показывается ошибка "Failed to load image"
- [ ] Сохранить оверлей → успешное сообщение
- [ ] Перезагрузить страницу → IMAGE слой загружается с правильным URL
- [ ] Изменить URL → обновляется в реальном времени
- [ ] Открыть оверлей на странице `/overlays` → изображение отображается

## Удаление отладочных логов

После того как проблема решена, удалите `console.log` из файлов:

1. `frontend/dashboard/src/features/overlay-builder/OverlayBuilder.vue`
2. `frontend/dashboard/src/features/overlay-builder/components/Canvas.vue`
3. `frontend/dashboard/src/features/overlay-builder/components/ImageLayerPreview.vue`
4. `frontend/dashboard/src/components/registry/overlays/edit.vue`

Найдите все `console.log` и `console.error` и удалите строки с отладкой.

## Что делать, если проблема не решена

1. Скопируйте логи из консоли браузера (все сообщения с префиксами [edit.vue], [OverlayBuilder], [Canvas], [ImageLayerPreview])
2. Откройте DevTools → Network → фильтр XHR/Fetch
3. Найдите GraphQL запросы (обычно `/graphql`)
4. Проверьте запрос и ответ - есть ли там `imageUrl`
5. Проверьте логи бэкенда на наличие ошибок
6. Проверьте базу данных напрямую - сохраняется ли `image_url`

## Полезные советы

### Хорошие источники изображений с CORS:
- `https://via.placeholder.com/300x200` - простые placeholder'ы
- `https://picsum.photos/300/200` - случайные фото
- Ваш собственный сервер с настроенным CORS

### Проверка CORS в DevTools:
1. F12 → Network
2. Загрузите изображение
3. Если видите ошибку типа "CORS policy" - сервер блокирует загрузку
4. Используйте другой источник изображений

### Пример правильного URL:
```
https://via.placeholder.com/300x200.png
```

### Примеры неправильных URL:
- `example.png` (относительный путь)
- `C:\images\photo.jpg` (локальный путь)
- `http://localhost/image.png` (только для разработки)
