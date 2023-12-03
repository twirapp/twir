export const chunk = <T>(arr: T[], size: number): T[][] => {
  const result: T[][]= [];

  for(let i = 0; i < arr.length; i += size) {
    result.push(arr.slice(i, i + size));
  }

  return result;
};
