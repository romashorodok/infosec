<script lang="ts">
	import { createDialog, melt } from "@melt-ui/svelte";
	/** Internal helpers */
	import { fade } from "svelte/transition";

	export let users: Promise<any>;
	export let boardTitle: string;
	export let boardId: string;
	export let onSubmit: (userID: string) => void;

	const {
		elements: {
			trigger,
			overlay,
			content,
			title,
			description,
			close,
			portalled,
		},
		states: { open },
	} = createDialog({
		forceVisible: false,
	});

	let userID: string;

	$: console.log(userID);
</script>

<button
	use:melt={$trigger}
	class="inline-flex items-center justify-center rounded-xl bg-white px-4 py-3
  font-medium leading-none text-magnum-700 shadow hover:opacity-75"
>
	Add User
</button>

<div use:melt={$portalled}>
	{#if $open}
		<div
			use:melt={$overlay}
			class="fixed inset-0 z-50 bg-black/50"
			transition:fade={{ duration: 150 }}
		/>
		<div
			class="fixed left-[50%] top-[50%] z-50 max-h-[85vh] w-[90vw]
            max-w-[450px] translate-x-[-50%] translate-y-[-50%] rounded-xl bg-white
            p-6 shadow-lg"
			use:melt={$content}
		>
			<h2
				use:melt={$title}
				class="m-0 text-lg font-medium text-black"
			>
				Add user to the board: {boardTitle} ({boardId})
			</h2>
			<p
				use:melt={$description}
				class="mb-5 mt-2 leading-normal text-zinc-600"
			>
				<br />
			</p>

			<fieldset class="mb-4 flex items-center gap-5">
				<label
					class="w-[90px] text-right text-black"
					for="username"
				>
					Username
				</label>

				{#await users}
					...Loading
				{:then users}
					<select bind:value={userID}>
						{#if users instanceof Array}
							{#each users as user}
								<option
									value={user.id}
								>
									{user.username}
								</option>
							{/each}
						{/if}
					</select>
				{/await}
			</fieldset>
			<div class="mt-6 flex justify-end gap-4">
				<button
					use:melt={$close}
					class="inline-flex h-8 items-center justify-center rounded-sm
                    bg-zinc-100 px-4 font-medium leading-none text-zinc-600"
				>
					Cancel
				</button>
				<button
					on:click={() => onSubmit(userID)}
					use:melt={$close}
					class="inline-flex h-8 items-center justify-center rounded-sm
                    bg-magnum-100 px-4 font-medium leading-none text-magnum-900"
				>
					Add user
				</button>
			</div>
		</div>
	{/if}
</div>
