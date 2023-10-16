<script lang="ts">
	export let onSubmit: (boardTitle: string) => void;

	let editable: boolean = false;

	let boardTitle: string;

	function onEdit() {
		editable = true;
	}
</script>

{#if editable}
	<div class="px-2">
		<input bind:value={boardTitle} class="w-[150px]" />
		<button
			on:click={() => {
				onSubmit(boardTitle);
				editable = false;
				boardTitle = undefined;
			}}>Submit</button
		>
		<button on:click={() => (editable = false)}>Cancel</button>
	</div>
{/if}

{#if !editable}
	<!-- svelte-ignore a11y-no-static-element-interactions -->
	<!-- svelte-ignore a11y-click-events-have-key-events -->
	<div on:click={onEdit}>
		<p>Create</p>
	</div>
{/if}
