import { useEffect, useState } from 'preact/compat';

import { browserUnProtectedClient } from '../api/twirp-browser.js';

export default () => {
	const [error, setError] = useState('');

	useEffect(() => {
		const url = new URL(window.location.href);
		const code = url.searchParams.get('code');
		const err = url.searchParams.get('error');
		if (err) {
			setError(err);
		} else {
			browserUnProtectedClient.authPostCode({
				code,
			})
				.then(() => window.location.replace('/dashboard'))
				.catch(() => setError('internal error'));
		}
	}, []);

	return <>
	Login
	{error}
</>;
};
