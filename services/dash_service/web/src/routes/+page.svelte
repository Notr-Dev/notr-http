<script>
  import { dev } from "$app/environment"
  import { fetchPath } from "$lib"
</script>

<h2>Services</h2>

{#await fetchPath("/services")}
  <p>loading...</p>
{:then data}
  {@const dashPath = dev ? "" : data.find((s) => s.package_id == "dash").path}
  {#each data as service}
    <a href="{dashPath}/services/{service.package_id}">
      {service.package_id} - {service.name} - {service.is_initialized
        ? "initialized"
        : "not initialized"}
    </a>
  {/each}
{:catch error}
  <p>{error.message}</p>
{/await}

<style>
  a {
    display: block;
  }
</style>
