import { RequestErrorCode, createRequestError } from "$lib/errors.js";

export async function fetchWithTimeout(
    input: RequestInfo | URL,
    init: RequestInit = {},
    timeoutMs: number = 5000
): Promise<Response> {
    const controller = new AbortController();
    const abort = controller.abort.bind(controller);
    init.signal?.addEventListener("abort", abort, { once: true });
    const timer = setTimeout(abort, timeoutMs);

    try {
        return await fetch(input, {
            ...init,
            signal: controller.signal,
        });
    } catch (err) {
        if (err && (err as Error).name === "AbortError") {
            throw createRequestError(RequestErrorCode.REQUEST_TIMEOUT, `Request timed out after ${timeoutMs} ms`);
        }
        throw err;
    } finally {
        clearTimeout(timer);
        init.signal?.removeEventListener("abort", abort);
    }
}
