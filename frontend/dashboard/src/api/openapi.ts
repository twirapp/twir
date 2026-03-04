import { Api, HttpClient, type FullRequestParams, type HttpResponse } from '@twir/api/openapi';
import { toast } from 'vue-sonner';

// Create a wrapper that intercepts requests and handles errors
function createErrorHandlingHttpClient(config: ConstructorParameters<typeof HttpClient>[0]) {
	const client = new HttpClient(config);
	const originalRequest = client.request.bind(client);

	client.request = async <T = any, E = any>(
		params: FullRequestParams,
	): Promise<HttpResponse<T, E>> => {
		try {
			return await originalRequest<T, E>(params);
		} catch (error: any) {
			// Handle HTTP errors
			if (error && error.error) {
				const errorData = error.error;

				// Handle standard error format (ErrorModel)
				if (errorData.title && errorData.detail) {
					const message = errorData.title;
					const description = errorData.detail;

					// If there's an errors array, format them
					if (errorData.errors && Array.isArray(errorData.errors)) {
						const details = errorData.errors
							.map((err: any) => {
								const location = err.location ? `${err.location}: ` : '';
								return `${location}${err.message}`;
							})
							.join(', ');

						toast.error(message, {
							description: details || description,
							duration: 10000,
						});
					} else {
						toast.error(message, {
							description,
							duration: 8000,
						});
					}

					throw error;
				}

				// Handle other error formats
				if (typeof errorData === 'string') {
					toast.error('Error', {
						description: errorData,
						duration: 5000,
					});
					throw error;
				}
			}

			// Network errors or unexpected errors
			toast.error('Network error', {
				description: 'Failed to connect to the server',
				duration: 5000,
			});

			throw error;
		}
	};

	return client;
}

export const openApi = new Api(
	createErrorHandlingHttpClient({
		baseUrl: `${window.location.origin}/api`,
		baseApiParams: {
			credentials: 'include',
		},
	}),
);
