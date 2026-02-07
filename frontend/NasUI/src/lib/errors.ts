export const AppErrorCode = {
    REQUEST_TIMEOUT: "REQUEST_TIMEOUT",
    FETCH_POOLS_FAILED: "FETCH_POOLS_FAILED",
    FETCH_POOL_FAILED: "FETCH_POOL_FAILED",
    FETCH_DRIVES_FAILED: "FETCH_DRIVES_FAILED",
    FETCH_ADOPTED_DRIVES_FAILED: "FETCH_ADOPTED_DRIVES_FAILED",
    CREATE_POOL_FAILED: "CREATE_POOL_FAILED",
    BUILD_POOL_FAILED: "BUILD_POOL_FAILED",
    DELETE_POOL_FAILED: "DELETE_POOL_FAILED",
} as const;

export type AppErrorCode = (typeof AppErrorCode)[keyof typeof AppErrorCode];

export class AppError extends Error {
    code: AppErrorCode;
    status?: number;

    constructor(code: AppErrorCode, message: string, status?: number) {
        super(message);
        this.code = code;
        this.name = code;
        this.status = status;
    }
}

export function createAppError(code: AppErrorCode, message: string, status?: number): AppError {
    return new AppError(code, message, status);
}

export async function responseError(
    response: Response,
    code: AppErrorCode,
    fallbackMessage: string,
): Promise<AppError> {
    let message = fallbackMessage;
    try {
        const data = await response.clone().json();
        if (data && typeof data.error === "string" && data.error.trim() !== "") {
            message = data.error;
        }
    } catch {
        // Ignore parse issues and keep fallback message.
    }
    return createAppError(code, `${message} (${response.status})`, response.status);
}
