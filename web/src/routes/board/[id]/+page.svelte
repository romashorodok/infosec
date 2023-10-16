<script lang="ts">
	import { API_ROUTE } from "$lib";
	import AddUserDialog from "$lib/components/add-user-dialog.svelte";
	import RemoveUserDialog from "$lib/components/remove-user-dialog.svelte";
	import { accessToken } from "$lib/stores/auth";
	import type { Fetch } from "$lib/utils/fetch";
	import { refreshPage } from "$lib/utils/refresh-page";
	import type { PageData } from "./$types";

	export let data: PageData;

	let boardError: boolean;

	$: if (data?.board?.message) {
		const msg = data.board.message;
		console.log(msg);
		boardError = true;
	} else {
		boardError = false;
	}

	$: ({ fetch } = data);

	const userFetcher = async (fetch: Fetch): Promise<any> => {
		try {
			return await fetch(`${API_ROUTE}/users`, {
				headers: {
					Authorization: `Bearer ${$accessToken}`,
				},
			}).then((r) => r.json());
		} catch (err) {
			console.log(err);
			return null;
		}
	};

	$: users = userFetcher(fetch);

	const onAddUserSubmit = async (userID: string) => {
		fetch(`${API_ROUTE}/boards/${data.board.id}/participants`, {
			method: "POST",
			headers: {
				Authorization: `Bearer ${$accessToken}`,
				"Content-Type": "application/json",
			},
			body: JSON.stringify({
				user_id: userID,
			}),
		}).then(refreshPage);
	};

	const onRemoveUserSubmit = async (userID: string) => {
		fetch(
			`${API_ROUTE}/boards/${data.board.id}/participants?user_id=${userID}`,
			{
				method: "DELETE",
				headers: {
					Authorization: `Bearer ${$accessToken}`,
					"Content-Type": "application/json",
				},
			}
		).then(refreshPage);
	};
</script>

<div class="overflow-y-scroll p-4">
	{#if !boardError}
		<div class="flex my-4">
			<AddUserDialog
				{users}
				boardTitle={data.board.title}
				boardId={data.board.id}
				onSubmit={onAddUserSubmit}
			/>

			<RemoveUserDialog
				{users}
				boardTitle={data.board.title}
				boardId={data.board.id}
				onSubmit={onRemoveUserSubmit}
			/>
		</div>

		<div class="mt-4 p-8 overflow-y-scroll bg-white rounded">
			<pre class="h-[400px]">{JSON.stringify(
					data.board,
					null,
					1
				)}</pre>
		</div>
	{:else}
		<h1>404 Not found</h1>
	{/if}
</div>
