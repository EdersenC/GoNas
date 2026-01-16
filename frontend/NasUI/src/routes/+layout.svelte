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
        background-color: #2c3e50;
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

<nav>
    <a href="/">Home</a>
    <a href="/pools">Pools</a>
    <a href="/drives">Drives</a>
    <a href="/pool/create">Pool Creator</a>
</nav>

{@render children()}