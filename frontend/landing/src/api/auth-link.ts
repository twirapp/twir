import { gqlUrl } from '@/api/gql.ts';

export async function getAuthLink(redirectTo: string) {
	const request = await fetch(gqlUrl, {
		'headers': {
			'content-type': 'application/json',
		},
		'body': `{"operationName":"LoginLink","query":"query LoginLink($redirectTo: String!) {\\n  authLink(redirectTo: $redirectTo)\\n}","variables":{"redirectTo":"${redirectTo}"}}`,
		'method': 'POST',
	});

	const response = await request.json();

	if (!request.ok || response.errors) {
		console.log(response);
		throw new Error(response.errors.toString());
	}

	return response.data.authLink;
}
