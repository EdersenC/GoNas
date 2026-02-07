export const RequestErrorCode = {
    REQUEST_TIMEOUT: "REQUEST_TIMEOUT",
} as const;

export const PoolErrorCode = {
    FETCH_POOLS_FAILED: "FETCH_POOLS_FAILED",
    FETCH_POOL_FAILED: "FETCH_POOL_FAILED",
    CREATE_POOL_FAILED: "CREATE_POOL_FAILED",
    BUILD_POOL_FAILED: "BUILD_POOL_FAILED",
    DELETE_POOL_FAILED: "DELETE_POOL_FAILED",
} as const;

export const DriveErrorCode = {
    FETCH_DRIVES_FAILED: "FETCH_DRIVES_FAILED",
    FETCH_ADOPTED_DRIVES_FAILED: "FETCH_ADOPTED_DRIVES_FAILED",
} as const;

export type RequestErrorCode = (typeof RequestErrorCode)[keyof typeof RequestErrorCode];
export type PoolErrorCode = (typeof PoolErrorCode)[keyof typeof PoolErrorCode];
export type DriveErrorCode = (typeof DriveErrorCode)[keyof typeof DriveErrorCode];
export type AppErrorCode = RequestErrorCode | PoolErrorCode | DriveErrorCode;

type ErrorKind = "request" | "pool" | "drive";

export class AppError extends Error {
    code: AppErrorCode;
    kind: ErrorKind;
    status?: number;

    constructor(kind: ErrorKind, code: AppErrorCode, message: string, status?: number) {
        super(message);
        this.code = code;
        this.kind = kind;
        this.name = code;
        this.status = status;
    }
}

export class RequestError extends AppError {
    constructor(code: RequestErrorCode, message: string, status?: number) {
        super("request", code, message, status);
    }
}

export class PoolError extends AppError {
    constructor(code: PoolErrorCode, message: string, status?: number) {
        super("pool", code, message, status);
    }
}

export class DriveError extends AppError {
    constructor(code: DriveErrorCode, message: string, status?: number) {
        super("drive", code, message, status);
    }
}

export function createRequestError(code: RequestErrorCode, message: string, status?: number): RequestError {
    return new RequestError(code, message, status);
}

function createByCode(code: AppErrorCode, message: string, status?: number): AppError {
    if (Object.values(PoolErrorCode).includes(code as PoolErrorCode)) {
        return new PoolError(code as PoolErrorCode, message, status);
    }
    if (Object.values(DriveErrorCode).includes(code as DriveErrorCode)) {
        return new DriveError(code as DriveErrorCode, message, status);
    }
    return new RequestError(code as RequestErrorCode, message, status);
}

export async function responseError(
    response: Response,
    code: PoolErrorCode | DriveErrorCode,
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
    return createByCode(code, `${message} (${response.status})`, response.status);
}

export function toUserMessage(err: unknown): string {
    if (err instanceof AppError) {
        if (err.code === RequestErrorCode.REQUEST_TIMEOUT) {
            return "The request took too long. Please try again.";
        }
        if (err.code === PoolErrorCode.CREATE_POOL_FAILED) {
            const raw = err.message.toLowerCase();
            if (raw.includes("invalid raid name")) {
                return "Pool name is invalid. Use letters, numbers, dots, dashes, or underscores only, with no spaces.";
            }
            if (raw.includes("pool created, but build failed")) {
                return "Pool was created, but the build failed. Open the pool and retry build.";
            }
            return "Could not create the pool. Please check your inputs and try again.";
        }
        if (err.code === PoolErrorCode.BUILD_POOL_FAILED) {
            return "Could not build the pool. Please verify the pool configuration and try again.";
        }
        if (err.code === PoolErrorCode.DELETE_POOL_FAILED) {
            return "Could not delete the pool right now. Please try again.";
        }
    }
    if (err instanceof Error && err.message.trim() !== "") {
        return err.message;
    }
    return "Something went wrong. Please try again.";
}
