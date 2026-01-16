<script lang="ts">
    interface Props {
        degraded?: boolean;
        offline?: boolean;
    }

    let { degraded = false, offline = false }: Props = $props();

    let status = $derived(degraded ? 'degraded' : offline ? 'offline' : 'online');

    const statusClasses: Record<string, string> = {
        degraded: 'bg-red-500',
        offline: 'bg-orange-500',
        online: 'bg-green-500'
    };
    console.log(statusClasses[status]);
</script>

<span class="w-3 h-3 rounded-full inline-flex shrink-0 self-start {statusClasses[status]}"></span>

