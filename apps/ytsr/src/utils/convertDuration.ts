export const durationToMilliseconds = (duration: string) => {
  const timeParts = duration.split(':').map(part => parseInt(part)); // convert each part to an integer

  const numComponents = timeParts.length;

  // If the duration has no components, return 0
  if (numComponents === 0) {
    return 0;
  }

  // If the duration has only one component (seconds), treat it as such
  if (numComponents === 1) {
    return timeParts[0]! * 1000;
  }

  // If the duration has two components (minutes, seconds), add the two components in milliseconds
  if (numComponents === 2) {
    return (timeParts[0]! * 60 + timeParts[1]!) * 1000;
  }

  // If the duration has three components (hours, minutes, seconds), add all three components in milliseconds
  if (numComponents === 3) {
    return (timeParts[0]! * 60 * 60 + timeParts[1]! * 60 + timeParts[2]!) * 1000;
  }

  // If the duration has more than three components, ignore the extra components and calculate the duration in milliseconds
  const factor = [216000, 3600, 60, 1];
  return timeParts.slice(0, 3).reduce((acc, val, index) => {
    return acc + (val * factor.slice(index).reduce((acc, val) => acc * val, 1) * 1000);
  }, 0);
};