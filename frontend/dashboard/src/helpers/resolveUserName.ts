export function resolveUserName(userName: string, userDisplayName?: string): string {
	if (!userDisplayName) return userName;

	if (userName === userDisplayName.toLocaleLowerCase()) {
		return userDisplayName;
	}

	return `${userDisplayName} (${userName})`;
}
