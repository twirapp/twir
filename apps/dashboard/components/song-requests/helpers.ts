export const formatDuration = (seconds: number) => {
  const format = (val: number) => `0${Math.floor(val)}`.slice(-2);
  const minutes = (seconds % 3600) / 60;

  return [minutes, seconds % 60].map(format).join(':');
};

export function millisToMinutesAndSeconds(millis: number) {
  const minutes = Math.floor(millis / 60000);
  const seconds = ((millis % 60000) / 1000);
  return minutes + ':' + (seconds < 10 ? '0' : '') + seconds.toFixed(0);
}