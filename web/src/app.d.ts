import type { IdentityTokenPayload } from '$lib/stores/auth';

declare global {
	namespace App {
		// interface Error {}
		// interface Platform {}

		interface PageData {
		}

		interface Locals {
			identityPayload: IdentityTokenPayload | null

			getAccessToken: function(): Promise<String | null>
		}
	}
}
export {};
