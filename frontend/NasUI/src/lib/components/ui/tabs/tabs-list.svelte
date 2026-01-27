<script lang="ts">
	import { Tabs as TabsPrimitive } from "bits-ui";
	import { cn } from "$lib/utils.js";

	let { class: className, ...restProps }: TabsPrimitive.ListProps = $props();
</script>

<TabsPrimitive.List
	data-slot="tabs-list"
	class={cn(
		"inline-flex h-10 items-center justify-start gap-1 rounded-xl border border-brand/40 bg-surface/60 p-1 text-muted-foreground",
		className
	)}
	{...restProps}
/>
