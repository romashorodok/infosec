import { API_ROUTE } from "$lib";
import type { Fetch } from "$lib/utils/fetch";
import type { ServerLoad } from "@sveltejs/kit";

async function loadBoard(fetch: Fetch, id: string): Promise<any | null> {
	try {
		return await fetch(`${API_ROUTE}/boards/${id}`).then(r => r.json())
	} catch (e) {
		console.log(e)
		return null
	}
}

export const load: ServerLoad = ({ params, fetch }) => {
	const { id } = params

	return {
		board: loadBoard(fetch, id)
	}
} 
