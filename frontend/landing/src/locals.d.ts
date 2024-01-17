import type { Profile } from '@twir/api/messages/auth/auth';

declare global {
	namespace App {
		interface Locals {
			profile?: Profile,
			authLink: string
		}
	}
}

export {};
