export const chunk = <T>(arr: T[], size: number): Array<T[]> => {
  return Array.from({ length: Math.ceil(arr.length / size) }, (_: any, i: number) =>
    arr.slice(i * size, i * size + size),
  );
};
