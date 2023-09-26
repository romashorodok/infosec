<script lang="ts">
	import { accessToken, identity, logout } from "$lib/stores/auth";
	import { env } from "$env/dynamic/public";

	async function verify() {
		const { data } = await fetch(
			`${env.PUBLIC_API_HOST}/token-revocation:verify`,
			{
				method: "post",
				headers: {
					Authorization: `Bearer ${$accessToken}`,
				},
			}
		).then((r) => r.json());

		console.log(data);
	}
</script>

<pre>{JSON.stringify($identity, null, 2)}</pre>

<button on:click|preventDefault={verify}>Refresh token</button>

<button on:click={logout}>Logout</button>
