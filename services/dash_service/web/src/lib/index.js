import { dev } from "$app/environment"

export async function fetchPath(path = "/") {
  if (path[0] !== "/") {
    throw new Error("Path must start with a /")
  }
  const url = new URL(
    "/dash" + path,
    dev ? "http://localhost:8080" : window.location.origin
  )
  const response = await fetch(url)
  if (!response.ok) {
    throw new Error("Network response was not ok")
  }
  const json = await response.json()
  console.log(url, json)
  return json
}
