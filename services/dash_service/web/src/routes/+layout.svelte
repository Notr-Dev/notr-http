<script>
  import { fetchPath } from "$lib"
  import "./globals.css"
</script>

<nav>
  {#await fetchPath("/details")}
    <p>loading...</p>
  {:then data}
    <h1>{data.name}</h1>
    <p>v{data.version}</p>
    <p>{data.port}</p>
  {:catch error}
    <p>{error.message}</p>
  {/await}
</nav>

<section>
  <slot />
</section>

<style>
  nav {
    background-color: #f2f2f2;
    padding: 1em;
    margin-bottom: 1em;
    display: flex;
    justify-content: space-between;
    gap: 1rem;
  }
  section {
    padding: 1rem;
  }
  p {
    display: flex;
    align-items: center;
  }
</style>
