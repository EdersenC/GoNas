<script lang="ts">
    import { AspectRatio } from "$lib/components/ui/aspect-ratio/index.ts";
    import {List} from "$lib/components/ui/drive/index.ts"
    import {Header} from "$lib/components/ui/card/index.ts";

</script>

<div class="grid grid-cols-1 gap-y-6 sm:gap-y-10">
    <div class="flex items-center justify-center px-2 sm:px-40 py-0">
        <div class="w-full">
            <AspectRatio ratio={2} class="bg-muted overflow-hidden">
                <div class="h-full flex flex-col min-h-0">
                    <Header class="p-2 shrink-0 bg-secondary text-secondary-foreground">
                        Drives
                    </Header>
                    <div class="flex-1 overflow-y-auto min-h-0">
                        <List />
                    </div>
                </div>
            </AspectRatio>
        </div>
    </div>

    <div class="flex items-center justify-center px-2 sm:px-40 py-0">
        <div class="w-full">
            <AspectRatio ratio={2} class="bg-muted overflow-hidden">
                <div class="h-full flex flex-col min-h-0">
                    <div class="p-2 shrink-0">Pools</div>
                    <div class="flex-1 overflow-y-auto min-h-0">
                        <List />
                    </div>
                </div>
            </AspectRatio>
        </div>
    </div>
</div>
