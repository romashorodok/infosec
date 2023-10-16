<script lang="ts">
	import { accessToken, identity } from "$lib/stores/auth";
	import { env } from "$env/dynamic/public";

	async function verify() {
		await fetch(`${env.PUBLIC_API_HOST}/token-revocation:verify`, {
			method: "post",
			headers: {
				Authorization: `Bearer ${$accessToken}`,
			},
		})
			.then((r) => r.json())
			.then(console.log);
	}
</script>

<div class="h-full">
	<pre>{JSON.stringify($identity, null, 2)}</pre>

	<button on:click|preventDefault={verify}>Refresh token</button>
</div>
