<script lang="ts">
    import type {Pool} from "$lib/models/pool.ts";

    interface Props {
		pool: Pool;
	}

	let { pool }: Props = $props();
</script>

<div class="relative overflow-hidden rounded-2xl border border-brand/50 bg-surface/80 p-5 shadow-sm">
	<div
		class="absolute inset-0 opacity-60"
		style="background: radial-gradient(circle at top, rgba(59, 130, 246, 0.22), transparent 60%);"
	></div>
	<div class="relative">
		<p class="text-xs font-semibold uppercase tracking-[0.25em] text-panel-foreground/70">
			Pool Manager
		</p>
		<h3 class="text-2xl font-semibold tracking-tight text-panel-foreground">
			Manage {pool.name || "Pool"}
		</h3>
		<p class="text-sm text-panel-foreground/70">Configure storage pool actions and settings.</p>
	</div>
</div>
