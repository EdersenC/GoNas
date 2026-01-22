<script lang="ts">
    interface Props {
        percent: number;
    }

    let { percent }: Props = $props();

    let color = $derived(percent < 50 ? 'green' : percent < 80 ? 'orange' : 'red');
    let strokeColor = $derived(
        color === 'green'
            ? 'var(--color-success)'
            : color === 'orange'
              ? 'var(--color-warning)'
              : 'var(--color-danger)'
    );

    const radius = 40;
    const circumference = 2 * Math.PI * radius;
    let strokeDasharray = $derived(`${circumference * (percent / 100)} ${circumference}`);
</script>

<div class="size-circle">
    <svg width="100" height="100" viewBox="0 0 100 100">
        <!-- Background circle -->
        <circle
            cx="50"
            cy="50"
            r="{radius}"
            fill="none"
            stroke="var(--color-surface-border)"
            stroke-width="8"
        />
        <!-- Progress circle -->
        <circle
            cx="50"
            cy="50"
            r="{radius}"
            fill="none"
            stroke="{strokeColor}"
            stroke-width="8"
            stroke-dasharray="{strokeDasharray}"
            stroke-linecap="round"
            transform="rotate(-90 50 50)"
        />
        <!-- Percentage text -->
        <text
            x="50"
            y="55"
            text-anchor="middle"
            font-size="16"
            fill="var(--color-surface-foreground)"
            font-weight="bold"
        >
            {Math.round(percent)}%
        </text>
    </svg>
</div>

<style>
    .size-circle {
        display: flex;
        justify-content: center;
        align-items: center;
    }
</style>
