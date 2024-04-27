declare global {
	namespace App {
		interface Locals {
			profile?: {
				displayName: string,
				profileImageUrl: string,
			},
			authLink: string
		}
	}
}

export {};
