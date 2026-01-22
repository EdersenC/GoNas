<script lang="ts">
    interface Props {
        status?: string;
    }

    let { status = 'unknown' }: Props = $props();

    // Normalize the status string
    let normalizedStatus = $derived(status?.toLowerCase() || 'unknown');

    const statusClasses: Record<string, string> = {
        'online': 'bg-success',
        'active': 'bg-success',
        'healthy': 'bg-success',
        'degraded': 'bg-warning',
        'rebuilding': 'bg-warning',
        'syncing': 'bg-warning',
        'offline': 'bg-danger',
        'failed': 'bg-danger',
        'error': 'bg-danger',
        'unknown': 'bg-muted'
    };

    let statusClass = $derived(statusClasses[normalizedStatus] || statusClasses['unknown']);
</script>

<span class="w-3 h-3 rounded-full inline-flex shrink-0 self-center {statusClass}" title={status}></span>
