<script lang="ts">
    import type { HTMLButtonAttributes } from 'svelte/elements';
    import { cn } from '$lib/utils.js';

    interface Props extends HTMLButtonAttributes {
        class?: string;
        isSpinning?: boolean;
        onclick?: () => void | Promise<void>;
    }

    let {
        class: className,
        isSpinning = false,
        onclick,
        ...restProps
    }: Props = $props();

    let spinning = $state(false);

    async function handleClick() {
        if (spinning) return;

        spinning = true;

        if (onclick) {
            await onclick();
        }

        // Complete one full rotation (1 second)
        setTimeout(() => {
            spinning = false;
        }, 1000);
    }
</script>

<button
    class={cn(
        "inline-flex items-center justify-center rounded-md text-sm font-medium transition-colors",
        "hover:bg-zinc-700 hover:text-zinc-100",
        "disabled:pointer-events-none disabled:opacity-50",
        "h-9 w-9 p-0",
        className
    )}
    onclick={handleClick}
    disabled={spinning || isSpinning}
    {...restProps}
>
    <svg
        class={cn(
            "h-5 w-5 transition-transform duration-500",
            (spinning || isSpinning) && "spinning"
        )}
        xmlns="http://www.w3.org/2000/svg"
        viewBox="0 0 24 24"
        stroke="currentColor"
        stroke-width="2"
        stroke-linecap="round"
        stroke-linejoin="round"
    >
        <path d="M21.5 2v6h-6M2.5 22v-6h6M2 11.5a10 10 0 0 1 18.8-4.3M22 12.5a10 10 0 0 1-18.8 4.2"/>
    </svg>
</button>

<style>
    @keyframes spin-once {
        from {
            transform: rotate(0deg);
        }
        to {
            transform: rotate(360deg);
        }
    }

    :global(.spinning) {
        animation: spin-once 1s ease-in-out forwards;
    }
</style>

