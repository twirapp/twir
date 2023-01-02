export const formatDuration = (seconds: number) => {
  const format = (val: number) => `0${Math.floor(val)}`.slice(-2);
  const minutes = (seconds % 3600) / 60;

  return [minutes, seconds % 60].map(format).join(':');
};