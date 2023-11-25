export const getTimeDiffInMilliseconds = (minutes: number) => {
	const startDate = new Date();
	const endDate = new Date(startDate.getTime() + (minutes * 60 * 1000));

	const diff = endDate.getTime() - startDate.getTime();

	return diff;
};

const formatNumber = (n: number) => n.toString().padStart(2, '0');
export const millisecondsToTime = (ms: number) => {
	const milliseconds = ms % 1000;
	ms = (ms - milliseconds) / 1000;
	const seconds = ms % 60;
	ms = (ms - seconds) / 60;
	const minutes = ms % 60;
	const hours = (ms - minutes) / 60;


	return `${hours ? formatNumber(hours)+':' : ''}${formatNumber(minutes)}:${formatNumber(seconds)}`;
};

