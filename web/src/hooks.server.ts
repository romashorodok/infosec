import { type Handle, type HandleFetch, json } from '@sveltejs/kit';

import { sequence } from '@sveltejs/kit/hooks';
import { mapCookiesFromHeader } from '$lib/utils/cookie';
import { env } from '$env/dynamic/public';

export const _REFRESH_TOKEN = '_refresh_token' as const;

export const VERIFY_ROUTE = `${env.PUBLIC_API_HOST}/token-revocation:verify` as const;
export const REFRESH_TOKEN_ROUTE = `${env.PUBLIC_API_HOST}/access-token` as const;

export const handleFetch: HandleFetch = (async ({ request, fetch, event: { cookies } }) => {

	const refreshToken = cookies.get(_REFRESH_TOKEN);

	if (!refreshToken)
		return json({ message: "missing refresh token" }, { status: 401 });

	request.headers.set('Authorization', `Bearer ${refreshToken}`);


	return fetch(request)
})

const verifyTokenHook: Handle = async ({ resolve, event }): Promise<any> => {
	const { cookies } = event;

	const refreshToken = cookies.get(_REFRESH_TOKEN);

	if (!refreshToken)
		return resolve(event)


	try {
		const response = await event.fetch(VERIFY_ROUTE, {
			method: 'POST'
		});

		if (response.status === 500) {
			cookies.delete(_REFRESH_TOKEN)
			return await resolve(event);
		}

		event.locals.identityPayload = await response.json();

		console.log(`Render page for user identity: ${JSON.stringify(event.locals.identityPayload)}`)
	} catch (e) {
		console.log(e)
	}

	return resolve(event)
}

const setFreshAccessTokenFuncHook: Handle = async ({ resolve, event }): Promise<any> => {
	const { cookies, fetch } = event;

	const refreshToken = cookies.get(_REFRESH_TOKEN);

	if (!refreshToken)
		return resolve(event)

	event.locals.getAccessToken = async (): Promise<String | null> => {
		try {
			const response = await fetch(REFRESH_TOKEN_ROUTE, { method: 'PUT' })

			const serverCookies = response.headers.get('set-cookie') as string

			await mapCookiesFromHeader(cookies, serverCookies);

			const { access_token } = await response.json()

			return access_token;

		} catch (e) {
			return null
		}
	}

	return resolve(event);
}

export const handle: Handle = sequence(
	verifyTokenHook,
	setFreshAccessTokenFuncHook,
) satisfies Handle;
