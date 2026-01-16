<script lang="ts">
    interface Props {
        status?: string;
    }

    let { status = 'unknown' }: Props = $props();

    // Normalize the status string
    let normalizedStatus = $derived(status?.toLowerCase() || 'unknown');

    const statusClasses: Record<string, string> = {
        'online': 'bg-green-500',
        'active': 'bg-green-500',
        'healthy': 'bg-green-500',
        'degraded': 'bg-yellow-500',
        'rebuilding': 'bg-yellow-500',
        'syncing': 'bg-yellow-500',
        'offline': 'bg-red-500',
        'failed': 'bg-red-500',
        'error': 'bg-red-500',
        'unknown': 'bg-gray-500'
    };

    let statusClass = $derived(statusClasses[normalizedStatus] || statusClasses['unknown']);
</script>

<span class="w-3 h-3 rounded-full inline-flex shrink-0 self-center {statusClass}" title={status}></span>

