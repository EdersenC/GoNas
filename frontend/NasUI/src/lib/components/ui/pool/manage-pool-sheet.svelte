<script lang="ts">
    import {Button} from "$lib/components/ui/button/index.ts";
    import type {Pool} from "$lib/models/pool.ts";
    import PoolSheetHeader from "./pool-sheet-header.svelte";
    import PoolStatsTab from "./pool-stats-tab.svelte";
    import PoolSharesTab from "./pool-shares-tab.svelte";
    import PoolActionsTab from "./pool-actions-tab.svelte";

    interface Props {
		pool: Pool;
		triggerClass?: string;
	}

	let { pool, triggerClass = "" }: Props = $props();

	let manageOpen = $state(false);
	let actionsTabRef: PoolActionsTab | undefined = $state();

	const handleSheetOpen = () => {
		actionsTabRef?.resetState?.();
	};

	const handleActionSuccess = () => {
		manageOpen = false;
	};
</script>

<Sheet.Root bind:open={manageOpen}>
	<Sheet.Trigger asChild>
		<Button
			variant="outline"
			class={`w-full text-surface-foreground border-surface-border hover:bg-surface-muted ${triggerClass}`}
			onclick={handleSheetOpen}
		>
			Manage Pool
		</Button>
	</Sheet.Trigger>
	<Sheet.Content side="right" class="bg-canvas text-canvas-foreground w-full p-0 sm:max-w-lg">
		<Sheet.Header class="sr-only">
			<Sheet.Title>Manage {pool.name || "Pool"}</Sheet.Title>
			<Sheet.Description>Configure storage pool actions and settings.</Sheet.Description>
		</Sheet.Header>
		<div class="flex h-full flex-col gap-4 bg-canvas p-5 text-canvas-foreground">
			<PoolSheetHeader {pool} />

			<Tabs.Root value="stats" class="flex flex-1 flex-col gap-4">
				<Tabs.List>
					<Tabs.Trigger value="stats">Stats</Tabs.Trigger>
					<Tabs.Trigger value="shares">Shares</Tabs.Trigger>
					<Tabs.Trigger value="actions">Actions</Tabs.Trigger>
				</Tabs.List>

				<Tabs.Content value="stats">
					<PoolStatsTab {pool} />
				</Tabs.Content>

				<Tabs.Content value="shares">
					<PoolSharesTab {pool} />
				</Tabs.Content>

				<Tabs.Content value="actions">
					<PoolActionsTab bind:this={actionsTabRef} {pool} onSuccess={handleActionSuccess} />
				</Tabs.Content>
			</Tabs.Root>
		</div>
	</Sheet.Content>
</Sheet.Root>
