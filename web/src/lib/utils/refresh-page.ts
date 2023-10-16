import { invalidateAll } from "$app/navigation";

export function refreshPage() {
	invalidateAll();
}
