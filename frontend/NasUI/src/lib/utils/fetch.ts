export async function fetchWithTimeout(
	input: RequestInfo | URL,
	init: RequestInit = {},
	timeoutMs: number = 5000
): Promise<Response> {
	const controller = new AbortController();
	const timer = setTimeout(() => controller.abort(), timeoutMs);

	if (init.signal) {
		if (init.signal.aborted) {
			controller.abort();
		} else {
			init.signal.addEventListener("abort", () => controller.abort(), { once: true });
		}
	}

	try {
		return await fetch(input, {
			...init,
			signal: controller.signal,
		});
	} catch (err) {
		if (err && (err as Error).name === "AbortError") {
			throw new Error(`Request timed out after ${timeoutMs} ms`);
		}
		throw err;
	} finally {
		clearTimeout(timer);
	}
}
