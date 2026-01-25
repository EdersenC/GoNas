<script lang="ts">
    import {UIDrive} from "$lib/components/ui/drive/index.js";
    import {onMount} from 'svelte';
    import { scale } from 'svelte/transition';
    import List from './list.svelte';
    import {DriveManagerKey, getDriveManagerContext, DriveManager} from "$lib/state/pool.svelte.js";

    // Cache drives data to prevent refetching on every mount

    let {
        ratio,
        poolCreatorMode,
    } = $props();

    let manager: DriveManager = getDriveManagerContext(DriveManagerKey);
    let adopted= true
    let error = $state<string | null>(null);
    let loading = $derived<boolean>(adopted ? manager.loadingAdoptedDrives : manager.loadingSystemDrives);

    onMount(async () => {
        await loadDrives()

    });


    async function loadDrives() {
        try {

            error = null;
            await manager.fetchAdoptedDrives()
            const count = 8;
            for (let i = 1; i <= count; i++) {
                const id = `fake-${i}`;
                let fakeDrive:Drive = {
                    name: `sd${i}`,
                    uuid: `fake-uuid-${i}`,
                    drive_key: { kind: 'by-path', value: `/dev/sd${i}` },
                    by_ids: [`/dev/disk/by-id/fake-${i}`],
                    wwid: `wwid-${i}`,
                    path: `/dev/sd${i}`,
                    size_sectors: 976773168,
                    logical_block_size: 512,
                    physical_block_size: 4096,
                    size_bytes: 500 * 1024 * 1024 * 1024*(i+2),
                    is_rotational: true,
                    model: `FAKE_MODEL_${i}`,
                    vendor: `FAKE_VENDOR`,
                    serial: `SNFAKE${1000 + i}`,
                   type: (i % 2 === 0 ? 'SSD' : 'HDD'),
                    mountpoint: '',
                    partitions: [],
                    fstype: '',
                    fsavail: (500 * 1024 * 1024 * 1024)*i,
                };
                manager.addAdoptedDrive(id, fakeDrive);
            }

        } catch (e) {
            console.error('Failed to load drives', e);
            error = 'Failed to load drives';
        }
    }

</script>

<List label="Drives" {loading} {error} onRefresh={loadDrives} ratio={ratio} >
    {#each Object.entries(manager.adoptedDrives) as [id, drive], i (id)}
        <div in:scale={{ duration: 300, delay: i * 50, start: 0.8 }}>
            <UIDrive
                    adopted={adopted}
                    drive={drive}
                    id={i}
                    poolCreatorMode={poolCreatorMode}
            />
        </div>
    {/each}
</List>