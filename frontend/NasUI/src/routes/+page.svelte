<script lang="ts">
    import {List, DrivesList} from "$lib/components/ui/drive/index.ts"
    import {PoolsList} from "$lib/components/ui/pool/index.ts"
    import {Separator} from "$lib/components/ui/separator/index.js";
</script>

<div class="grid grid-cols-1 gap-10 pt-10 pb-20">
    <DrivesList />
    <Separator class="mx-auto w-[calc(100%-4rem)] sm:w-[calc(100%-12rem)]" />
    <PoolsList />

</div>
