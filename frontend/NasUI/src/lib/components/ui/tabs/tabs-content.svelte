<script lang="ts">
	import { Tabs as TabsPrimitive } from "bits-ui";
	import { cn } from "$lib/utils.js";

	let { class: className, ...restProps }: TabsPrimitive.ContentProps = $props();
</script>

<TabsPrimitive.Content
	data-slot="tabs-content"
	class={cn("focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-brand/40", className)}
	{...restProps}
/>
