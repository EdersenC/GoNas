<script lang="ts">
	import { Separator as SeparatorPrimitive } from "bits-ui";
	import { cn } from "$lib/utils.js";

	let {
		ref = $bindable(null),
		class: className,
		"data-slot": dataSlot = "separator",
		...restProps
	}: SeparatorPrimitive.RootProps = $props();
</script>

<SeparatorPrimitive.Root
	bind:ref
	data-slot={dataSlot}
	class={cn(
		"bg-border/60 shrink-0 data-[orientation=horizontal]:h-[2px] data-[orientation=vertical]:min-h-full data-[orientation=vertical]:w-[2px]",
		className
	)}
	{...restProps}
/>
