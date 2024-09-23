import { error } from "@sveltejs/kit"

export async function load({ parent, params }) {
  const id = params.id
  return {
    id: id,
  }
}

export const prerender = false
