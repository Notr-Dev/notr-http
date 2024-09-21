export async function fetchPath(path = "/") {
  if (path[0] !== "/") {
    throw new Error("Path must start with a /")
  }
  const url = new URL("/dash" + path, window.location.origin)
  const response = await fetch(url)
  if (!response.ok) {
    throw new Error("Network response was not ok")
  }
  return await response.json()
}
