<script lang="ts">
    import {Button} from "$lib/components/ui/button/index.ts";
    import {Status, Size, ActionDropdown} from "$lib/components/ui/drive/index.ts";
    import { Root as CardRoot, Header as CardHeader, Content as CardContent, Footer as CardFooter, Title as CardTitle } from "$lib/components/ui/card/index.ts";
    let{
        drive,
        id,
        poolCreatorMode,
    } = $props();

   let realstate = $state(poolCreatorMode);

    function formatBytes(bytes: number): string {
        if (bytes === 0) return '0 B';
        const k = 1024;
        const sizes = ['B', 'KB', 'MB', 'GB', 'TB'];
        const i = Math.floor(Math.log(bytes) / Math.log(k));
        return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i];
    }

   export async function Post(driveId: string) {
        if (!driveId) return;
        // choose endpoint based on location
        let url = `http://localhost:8080/api/v1/drives/adopt/${driveId}`;
        if (window.location.pathname.startsWith('/pools')) {
            url = `http://localhost:8080/api/v1/pools/adopt/${driveId}`;
        }
        const res = await fetch(url, {
            method: 'POST',
        });
        console.log(await res.json())
        console.log(`Adopting drive with ID: ${driveId}`);
    }

    // Reactive computed values using the project's $derived helper (single-argument form)
    let used = $derived(drive ? drive.size_bytes - drive.fsavail : 0);
    let percent = $derived(drive && drive.size_bytes > 0 ? (used / drive.size_bytes) * 100 : 0);

</script>

{#if drive}
    <div class="drive-card w-full max-w-full min-w-0" id={"drive-" + id}>
        <CardRoot
                class="h-full w-full max-w-full min-w-0 flex flex-col
           !bg-panel !text-panel-foreground text-panel-foreground
           border border-panel-border/60 rounded-lg shadow-sm
           transform transition-transform transition-shadow transition-colors
           duration-100 ease-out will-change-transform
           hover:scale-[1.02] hover:shadow-lg hover:border-brand/50 hover:ring-2 hover:ring-brand/30{(poolCreatorMode) ? ' ring-2 ring-brand/50' : ''}"
                style="--card: var(--color-panel); --card-foreground: var(--color-panel-foreground); --card-border: var(--color-panel-border);"
                tabindex="0"
                role="button"
                onclick={()=> {console.log("Drive card clicked: " + (drive?.name || "unknown"))}}
                onkeydown={(e) => { if (e.key === 'Enter' || e.key === ' ') { e.preventDefault(); } }}
        >
            <CardHeader class="min-w-0">
                <!-- min-w-0 on the flex row is critical -->
                <CardTitle class="flex items-center gap-2 min-w-0">
                    <div class="flex items-center gap-2">
                        <!-- checkbox removed: entire card is clickable now -->
                    </div>
                    {#if drive.name?.includes("loop")}
                        <Status degraded={false} offline={false} />
                    {:else}
                        <Status degraded={false} offline={drive.is_rotational} />
                    {/if}

                    <!-- min-w-0 + truncate on text column -->
                    <span class="flex-1 min-w-0 text-center font-bold truncate">
          {drive.name}
        </span>

                    <!-- prevent dropdown from forcing width -->
                    <div class="shrink-0" on:click|stopPropagation>
                        <ActionDropdown />
                    </div>
                </CardTitle>
            </CardHeader>

            <CardContent class="flex-1 flex flex-col gap-1 text-sm min-w-0">
                <div class="flex justify-center mb-1 min-w-0">
                    <div class="flex flex-col gap-0 text-center min-w-0">
                        <span class="text-surface-foreground font-bold underline text-xs">Size</span>
                        <span class="font-bold text-sm truncate">
            {formatBytes(drive?.size_bytes || 0)}
          </span>
                    </div>
                </div>

                <!-- add min-w-0 to this row so children can shrink -->
                <div class="flex items-center min-w-0">
                    <div class="flex-1 text-left min-w-0">
                        <div class="flex flex-col gap-1 min-w-0">
                            <span class="text-surface-foreground font-bold underline text-xs">Used</span>
                            <span class="font-bold truncate">
              {formatBytes(used)}
            </span>
                        </div>
                    </div>

                    <div class="shrink-0 flex justify-center">
                        <Size percent={percent} />
                    </div>

                    <div class="flex-1 text-center min-w-0">
                        <div class="flex flex-col gap-1 min-w-0">
                            <span class="text-surface-foreground font-bold underline text-xs">Available</span>
                            <span class="font-bold truncate">
              {formatBytes(drive?.fsavail || 0)}
            </span>
                        </div>
                    </div>
                </div>
            </CardContent>

            <CardFooter class="min-w-0">
                <!-- stop long labels from expanding the card -->
                <div on:click|stopPropagation class="w-full">
                    <Button
                            variant="green"
                            class="w-full min-w-0 truncate"
                            onclick={() => Post((drive?.drive_key?.kind ?? "") + ":" + (drive?.drive_key?.value ?? ""))}
                            title={"Adopt " + (drive?.name ?? "")}
                    >
                        Adopt {drive.name}
                    </Button>
                </div>
            </CardFooter>
        </CardRoot>
    </div>

{:else}
    <div class="drive-card w-full max-w-full min-w-0" id={"drive-" + (id || "loading")}>
        <CardRoot class="h-full w-full max-w-full min-w-0 flex flex-col !bg-panel !text-panel-foreground text-panel-foreground border border-panel-border/60 rounded-lg shadow-sm">
            <CardContent class="flex-1 flex justify-center items-center min-w-0">
                <span class="truncate">Loading...</span>
            </CardContent>
        </CardRoot>
    </div>
{/if}


<style>
    .drive-card {
        width: 280px;
        height: 360px;
        display: block;
    }
</style>
