# Quick Fix Instructions - IMAGE Layer Display Issue

## Что было сделано

Я исправил проблему с отображением IMAGE слоёв и добавил подробное логирование для отладки.

## Изменённые файлы

### Критические исправления:

1. **frontend/dashboard/src/features/overlay-builder/components/ImageLayerPreview.vue**
   - Добавлена валидация `hasValidUrl` - теперь не пытается загрузить пустую строку как изображение
   - Добавлено логирование для отладки

2. **frontend/overlays/src/components/image-layer.vue**
   - Добавлена валидация `hasValidUrl` для рендера оверлея
   - Предотвращает попытку загрузки пустых URL

3. **frontend/dashboard/src/features/overlay-builder/components/PropertiesPanel.vue**
   - Исправлен обработчик события: `@update="emit('update', $event)"` вместо несуществующей функции

### Отладочное логирование (временно):

4. **frontend/dashboard/src/features/overlay-builder/OverlayBuilder.vue**
5. **frontend/dashboard/src/features/overlay-builder/components/Canvas.vue**
6. **frontend/dashboard/src/components/registry/overlays/edit.vue**

## Что делать дальше

### Шаг 1: Проверьте работу

1. Откройте дашборд
2. Откройте консоль браузера (F12)
3. Создайте или откройте оверлей с IMAGE слоем
4. Посмотрите на логи в консоли с префиксами:
   - `[edit.vue]` - загрузка/сохранение через API
   - `[OverlayBuilder]` - загрузка в редактор
   - `[Canvas]` - рендер слоя
   - `[ImageLayerPreview]` - компонент изображения

### Шаг 2: Диагностика по логам

Найдите в консоли строки типа:
```
[ImageLayerPreview] Component mounted with imageUrl: '...' hasValidUrl: true/false
```

**Если `hasValidUrl: false`** и `imageUrl` пустой:
- Это нормально для нового слоя
- Введите URL в панели Properties → Image URL
- Или нажмите "Use Placeholder Image"

**Если `hasValidUrl: true`** но показывается "Failed to load image":
- URL неверный или изображение не существует
- Или есть CORS блокировка
- Попробуйте: `https://via.placeholder.com/300x200`

**Если `imageUrl: undefined`**:
- Проблема на бэкенде - GraphQL не возвращает поле `imageUrl`
- Нужно проверить резолвер

### Шаг 3: Проверьте сохранение

1. Установите URL изображения в редакторе
2. Сохраните оверлей
3. В консоли найдите:
   ```
   [edit.vue] Layers to be saved: [...]
   ```
4. Убедитесь, что `imageUrl` присутствует в сохраняемых данных

5. Перезагрузите страницу
6. Проверьте, что изображение загрузилось обратно

### Шаг 4: Проверьте отображение в оверлее

1. Откройте страницу оверлея (не дашборд)
2. Проверьте, отображается ли изображение
3. Если нет - откройте консоль и проверьте ошибки

## Возможные проблемы и решения

### Проблема: "No image URL" в редакторе
**Решение**: Установите URL в панели Properties

### Проблема: "Failed to load image"
**Причины**:
- Неверный URL
- CORS блокировка
- Изображение не существует

**Решения**:
1. Проверьте URL - откройте его в новой вкладке браузера
2. Используйте CORS-friendly хостинги:
   - `https://via.placeholder.com/300x200`
   - `https://picsum.photos/300/200`
3. DevTools → Network → посмотрите на ошибку

### Проблема: imageUrl не сохраняется
**Проверка**:
1. Консоль браузера → найдите `[edit.vue] Layers to be saved`
2. Если `imageUrl` там есть - проблема на бэкенде
3. Если нет - проблема в редакторе

**Решение**:
- Проверьте, что используете компонент `ImageLayerEditor` для редактирования
- Проверьте, что событие `@update` правильно обрабатывается

## Удаление отладочных логов

После того как всё заработает, удалите `console.log` из файлов:

```bash
# Найти все отладочные логи:
grep -n "console.log.*\[edit.vue\]\|\[OverlayBuilder\]\|\[Canvas\]\|\[ImageLayerPreview\]" \
  frontend/dashboard/src/features/overlay-builder/OverlayBuilder.vue \
  frontend/dashboard/src/features/overlay-builder/components/Canvas.vue \
  frontend/dashboard/src/features/overlay-builder/components/ImageLayerPreview.vue \
  frontend/dashboard/src/components/registry/overlays/edit.vue
```

Удалите строки с этими логами вручную.

## Дополнительная документация

Создал 3 документа:

1. **IMAGE_LAYER_DEBUG_GUIDE.md** - Подробное руководство по отладке (английский)
2. **IMAGE_LAYER_FIX_RU.md** - Руководство по исправлению проблем (русский)
3. **IMAGE_LAYER_BUG_FIX_SUMMARY.md** - Полное описание всех изменений

## Быстрый тест

```bash
# 1. Пересобрать фронтенд
cd frontend/dashboard
bun install
bun run build

# 2. Открыть дашборд
# 3. Создать IMAGE слой
# 4. Установить URL: https://via.placeholder.com/300x200
# 5. Сохранить
# 6. Перезагрузить страницу
# 7. Проверить что изображение загрузилось
```

## Если проблема не решена

Пришлите мне:
1. Логи из консоли браузера (все с префиксами [edit.vue], [OverlayBuilder], [Canvas], [ImageLayerPreview])
2. Скриншот DevTools → Network → вкладка с GraphQL запросами
3. Описание того, что видите на экране

Я помогу разобраться дальше!

## Контакты для вопросов

Если что-то непонятно или нужна помощь - пишите в том же чате, приложу логи и помогу разобраться.
