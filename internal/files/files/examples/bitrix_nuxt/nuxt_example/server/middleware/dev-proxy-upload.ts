export default defineEventHandler(async (event) => {
  if (process.env.NODE_ENV !== 'development') return

  const path = event.path

  if (!/^\/(upload|api|local|bitrix)(\/|$)/.test(path)) return

  const target = `http://nginx${path}`
  const response = await proxyRequest(event, target, { fetch })
  return sendProxy(event, response)
})
