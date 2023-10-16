import { API_ROUTE } from "$lib";
import type { Fetch } from "$lib/utils/fetch";
import type { LayoutServerLoad } from "./$types";

async function loadUserBoards(fetch: Fetch) {
	try {
		return fetch(`${API_ROUTE}/users/boards`).then(r => r.json())
	}
	catch (_) {
		return null
	}
}

export const load: LayoutServerLoad = async ({ fetch, locals }) => {
	return {
		identity: locals.identityPayload || null,
		accessToken: locals.getAccessToken ? await locals.getAccessToken() : null,
		boards: loadUserBoards(fetch),
	}
}
