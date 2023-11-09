
export default {
  async email(message, env, ctx) {
    const timestamp = Date.now().toString(16);
    const random = Math.random().toString(16).substring(2);
    const randid = `D1-${timestamp}-${random}`;

    const reader = message.raw.getReader();
    const decoder = new TextDecoder();
    let result = '';
    while (true) {
      const { done, value } = await reader.read();
      if (done) break;
      result += decoder.decode(value, { stream: !reader.closed });
    }
  
    await env.ekv.prepare('INSERT INTO mail (id, fr1, to1, raw) VALUES (?1, ?2, ?3, ?4)')
                  .bind( randid, message.from, message.to, result ).run()
  },
  async fetch(request, env) {
    const url = new URL(request.url);
    if (url.pathname.endsWith("/eml")) {
      let ak = url.searchParams.get("ak")
      if (ak !== "api_xxx") {
        return new Response("token is invalid", { status: 403 });
      }
      let size = url.searchParams.get("size")
      if (size === undefined || size === null) { size = "10" }
      let rs = await env.ekv.prepare('SELECT id, raw FROM mail LIMIT ' + size).all()
      let rm = url.searchParams.get("rm")
      if (rm === "1") {
        for (const value of rs.results ) {
          await env.ekv.prepare('DELETE FROM mail WHERE id = ?1').bind(value.id).run()
        }
      }
      const data = {
        success: true,
        data: rs.results
      }
      const json = JSON.stringify(data, null, 2);
      return new Response(json, {
        headers: {
          "content-type": "application/json;charset=UTF-8",
        },
        status: 200
      });
    }
    return new Response("function not found", { status: 404 });
  }
}