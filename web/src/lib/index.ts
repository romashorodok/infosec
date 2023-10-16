// place files you want to import through the `$lib` alias in this folder.
// 

import { env } from "$env/dynamic/public";

export const REFRESH_TOKEN_ROUTE = `${env.PUBLIC_API_HOST}/access-token` as const;
export const API_ROUTE = `${env.PUBLIC_API_HOST}` as const;
