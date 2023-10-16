<script lang="ts">
	import "../app.css";

	import { accessToken, logout } from "$lib/stores/auth";
	import type { PageData } from "./$types";
	import BoardCreateCard from "$lib/components/board-create-card.svelte";
	import { API_ROUTE } from "$lib";
	import { refreshPage } from "$lib/utils/refresh-page";

	type Board = {
		id: number;
		title: string;
	};

	export let data: PageData;
	$: if (data?.boards && !data.boards.message) {
		boards = data.boards;
	}

	let boards: Array<Board> | null = null;

	$: ({ fetch } = data);

	$: console.log(boards)
	const createBoard = async (boardTitle: string) => {
		await fetch(`${API_ROUTE}/boards`, {
			method: "post",
			headers: {
				Authorization: `Bearer ${$accessToken}`,
				"Content-Type": "application/json",
			},
			body: JSON.stringify({
				title: boardTitle,
			}),
		}).then(refreshPage);
	};
</script>

<div class="contents min-h-inherit">
	<div class="flex flex-col min-h-inherit min-h-screen m-2">
		<nav
			class="flex flex-1 p-3 max-h-[50px] min-h-[50px] justify-between items-center z-[100] rounded-md"
		>
			<div>
				<a href="/">Home</a>
			</div>
			<div>
				<a href="/login">Login</a>

				<button on:click={logout}>Logout</button>
			</div>
		</nav>
		<div
			class="flex flex-row flex-1 pt-3 box-border overflow-hidden"
		>
			{#if boards}
				<aside class="w-[220px]">
					<BoardCreateCard
						onSubmit={createBoard}
					/>
					{#each boards as board}
						<div>
							<a
								href={`/board/${board.id}`}
							>
								{board.title}
							</a>
						</div>
					{/each}
				</aside>
			{/if}

			<main class="flex-1 overflow-y-scroll bg-gray-100">
				<slot />
			</main>
		</div>
	</div>
</div>

<style lang="postcss">
	:global(html) {
		min-height: 100%;
	}
	:global(body) {
		min-height: 100vh;
		overflow: hidden;
	}
</style>
