import { intervalToDuration, format, formatDuration } from 'date-fns';

export function humanizeStreamDuration(start: number, locale = 'en') {
  const duration = intervalToDuration({ start: start, end: Date.now() });

  return formatDuration(duration);
}
