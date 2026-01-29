<script lang="ts">
    import type {Pool} from "$lib/models/pool.ts";

    interface Props {
		pool: Pool;
	}

	let { pool }: Props = $props();
</script>

<div class="flex flex-col gap-4">
	<div class="rounded-2xl border border-brand/40 bg-surface/60 p-4 text-sm">
		<span class="text-xs font-semibold uppercase tracking-wider text-panel-foreground/70">
			Shares
		</span>
		<p class="mt-2 text-sm text-panel-foreground/70">
			Create shares within this pool for VM access, permissions, and quotas.
		</p>
		<div
			class="mt-3 rounded-lg border border-dashed border-brand/40 bg-surface-muted/30 px-3 py-4 text-xs text-panel-foreground/70"
		>
			Share management UI coming next.
		</div>
	</div>
</div>
