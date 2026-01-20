<script lang="ts">
    import * as DropdownMenu from "$lib/components/ui/dropdown-menu/index.js";
    import { Button } from "$lib/components/ui/button/index.js";

    let position = $state("bottom");
</script>

<DropdownMenu.Root>
    <DropdownMenu.Trigger asChild>
        <Button variant="ghost" size="sm" class="h-8 w-8 p-0">
            <span class="text-lg">⋮</span>
        </Button>
    </DropdownMenu.Trigger>
    <DropdownMenu.Content class="w-56">
        <DropdownMenu.Group>
            <DropdownMenu.Separator />
            <DropdownMenu.Item class="cursor-pointer" on:select={() => console.log("Action 1 selected")}>
                Action 1
            </DropdownMenu.Item>
        </DropdownMenu.Group>
    </DropdownMenu.Content>
</DropdownMenu.Root>
