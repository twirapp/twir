import type { Profile } from '@twir/grpc/generated/api/api/auth';

declare global {
	namespace App {
		interface Locals {
			profile?: Profile,
			authLink: string
		}
	}
}

export {};
