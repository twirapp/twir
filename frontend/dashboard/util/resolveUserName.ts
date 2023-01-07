export const resolveUserName = (name: string, displayName: string | null | undefined) => {
  if (!displayName) return name;

  if (name === displayName.toLowerCase()) return displayName;
  else return name;
};