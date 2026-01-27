<script lang="ts">
	import { Tabs as TabsPrimitive } from "bits-ui";
	import { cn } from "$lib/utils.js";

	let { class: className, ...restProps }: TabsPrimitive.TriggerProps = $props();
</script>

<TabsPrimitive.Trigger
	data-slot="tabs-trigger"
	class={cn(
		"inline-flex items-center justify-center rounded-lg px-3 py-2 text-xs font-semibold uppercase tracking-wider transition focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-brand/40 hover:bg-surface-muted/70 hover:text-surface-foreground hover:shadow-sm active:shadow data-[state=active]:border data-[state=active]:border-brand/60 data-[state=active]:bg-brand/15 data-[state=active]:text-panel-foreground data-[state=active]:shadow-md",
		className
	)}
	{...restProps}
/>
