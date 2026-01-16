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
		"bg-border/100 shrink-0 data-[orientation=horizontal]:h-[2px] data-[orientation=vertical]:min-h-full data-[orientation=vertical]:w-[0.04px]",
		className
	)}
	{...restProps}
/>
