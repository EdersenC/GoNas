<script lang="ts">
    let { children } = $props();
    import '../app.css'
</script>

<style>
    :global(body) {
        background-color: #121212;
        color: #9e0000;
    }

    nav {
        background-color: #2b7efc;
        padding: 1rem;
        display: flex;
        gap: 1rem;
    }

    nav a {
        color: #ffffff;
        text-decoration: none;
        font-weight: bold;
    }

    nav a:hover {
        text-decoration: underline;
    }


</style>


{@render children()}