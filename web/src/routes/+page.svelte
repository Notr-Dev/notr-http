<script>
  import { fetchPath } from "$lib"
</script>

{#await fetchPath("/details")}
  <p>loading...</p>
{:then data}
  {JSON.stringify(data)}
{:catch error}
  <p>{error.message}</p>
{/await}
