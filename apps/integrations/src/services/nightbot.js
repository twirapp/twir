export class Nightbot {

	accessToken;

	/**
	 *
	 * @param {string} accessToken
	 */
	constructor(accessToken) {
		this.accessToken = accessToken;
	}

	async getCustomCommands() {
		const req = await fetch('https://api.nightbot.tv/1/commands', {
      method: 'GET',
      headers: {
        'Authorization': `Bearer ${this.accessToken}`,
        },
    });

		const data = await req.json();

		if (!data.commands) {
			throw new Error('incorrect response');
		}

		return data.commands;
	}

	async destroy() {

	}
}
