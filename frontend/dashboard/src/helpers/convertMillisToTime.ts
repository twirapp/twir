export const padTo2Digits = (num: number) => {
	return num.toString().padStart(2, '0');
};

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

const rtf = new Intl.RelativeTimeFormat(window.navigator.language, {
	localeMatcher: 'best fit',
	numeric: 'always',
	style: 'long',
});

export const timeAgo = (value: string) => {
  const seconds = Math.floor((new Date().getTime() - new Date(value).getTime()) / 1000);
  let interval = seconds / 31536000;
  if (interval > 1) {
	return rtf.format(-Math.floor(interval), 'year');
  }
  interval = seconds / 2592000;
  if (interval > 1) {
	return rtf.format(-Math.floor(interval), 'month');
  }
  interval = seconds / 86400;
  if (interval > 1) {
	return rtf.format(-Math.floor(interval), 'day');
  }
  interval = seconds / 3600;
  if (interval > 1) {
	return rtf.format(-Math.floor(interval), 'hour');
  }
  interval = seconds / 60;
  if (interval > 1) {
	return rtf.format(-Math.floor(interval), 'minute');
  }

  return rtf.format(-Math.floor(interval), 'second');
};
