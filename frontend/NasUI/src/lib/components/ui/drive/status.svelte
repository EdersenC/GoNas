<script lang="ts">
    interface Props {
        degraded?: boolean;
        offline?: boolean;
    }

    let { degraded = false, offline = false }: Props = $props();

    let status = $derived(degraded ? 'degraded' : offline ? 'offline' : 'online');

    const statusClasses: Record<string, string> = {
        degraded: 'bg-warning',
        offline: 'bg-danger',
        online: 'bg-success'
    };
    console.log(statusClasses[status]);
</script>

<span class="w-3 h-3 rounded-full inline-flex shrink-0 self-center {statusClasses[status]}"></span>
