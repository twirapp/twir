export const createUserName = (userName: string, userDisplayName?: string) => {
	if (!userDisplayName) return userName;

	if (userName === userDisplayName.toLocaleLowerCase()) {
		return userDisplayName;
	}

	return `${userDisplayName} (${userName})`;
};
