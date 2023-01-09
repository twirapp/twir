export const convertMillisToTime = (millis: number) => {
  let seconds = Math.floor(millis / 1000);
  let minutes = Math.floor(seconds / 60);
  const hours = Math.floor(minutes / 60);

  seconds = seconds % 60;
  minutes = minutes % 60;

  return `${hours ? `${padTo2Digits(hours)}:` : ''}${padTo2Digits(minutes)}:${padTo2Digits(
    seconds,
  )}`;
};

export const padTo2Digits = (num: number) => {
  return num.toString().padStart(2, '0');
};

export const moveItem = <T>(items: T[], from: number, to: number) => {
  const values = [...items];
  values.splice(to, 0, values.splice(from, 1)[0]);
  return values;
};

export const toFixedNum = (num: number) => {
  return parseInt(num.toFixed(0), 10);
};

export function createdAtTime(createdAt: string | Date) {
  const date = createdAt instanceof Date ? createdAt : new Date(createdAt);
  const formatter = new Intl.RelativeTimeFormat();
  const ranges = {
    years: 3600 * 24 * 365,
    months: 3600 * 24 * 30,
    weeks: 3600 * 24 * 7,
    days: 3600 * 24,
    hours: 3600,
    minutes: 60,
    seconds: 1,
  } as Record<string, number>;

  const secondsElapsed = (date.getTime() - Date.now()) / 1000;

  for (const range in ranges) {
    if (ranges[range] < Math.abs(secondsElapsed)) {
      const delta = secondsElapsed / ranges[range];
      // eslint-disable-next-line @typescript-eslint/ban-ts-comment
      // @ts-ignore
      return formatter.format(Math.round(delta), range);
    }
  }
}
