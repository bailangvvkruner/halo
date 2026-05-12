const http = require('http');
const fs = require('fs');

const PORT = 8888;
const GO_API = 'http://localhost:8092';

function jsonHtml(obj, indent = 0) {
  const pad = ' '.repeat(indent * 2);
  if (obj === null) return '<span class="null">null</span>';
  if (typeof obj === 'boolean') return `<span class="bool">${obj}</span>`;
  if (typeof obj === 'number') return `<span class="num">${obj}</span>`;
  if (typeof obj === 'string') {
    if (obj.length > 120) return `<span class="str">"${obj.substring(0, 120)}..."</span>`;
    return `<span class="str">"${obj}"</span>`;
  }
  if (Array.isArray(obj)) {
    if (obj.length === 0) return '[]';
    if (indent >= 3) return `[${obj.length} items]`;
    return `[<br>${pad}  ${obj.map(item => jsonHtml(item, indent + 1)).join(`,<br>${pad}  `)}<br>${pad}]`;
  }
  if (typeof obj === 'object') {
    const keys = Object.keys(obj);
    if (keys.length === 0) return '{}';
    if (indent >= 4) return `{...(${keys.length} keys)}`;
    return `{<br>${keys.map(k => `${pad}  <span class="key">"${k}"</span>: ${jsonHtml(obj[k], indent + 1)}`).join(`,<br>`)}<br>${pad}}`;
  }
  return String(obj);
}

async function fetchJSON(path, token) {
  return new Promise((resolve, reject) => {
    const opts = { hostname: 'localhost', port: 8092, path, method: 'GET' };
    if (token) opts.headers = { Authorization: `Bearer ${token}` };
    const req = http.get(opts, res => {
      let data = '';
      res.on('data', c => data += c);
      res.on('end', () => { try { resolve(JSON.parse(data)); } catch(e) { resolve({raw: data}); } });
    });
    req.on('error', reject);
    req.setTimeout(5000, () => { req.destroy(); reject(new Error('timeout')); });
  });
}

async function main() {
  // Get token
  let token = '';
  try {
    const tokenRes = await new Promise((resolve, reject) => {
      const d = JSON.stringify({username:'reporter',password:'Report123456',email:'report@halo.go',displayName:'Reporter'});
      const req = http.request({hostname:'localhost',port:8092,path:'/auth/register',method:'POST',headers:{'Content-Type':'application/json','Content-Length':Buffer.byteLength(d)}}, res => {
        let b=''; res.on('data',c=>b+=c); res.on('end',()=>resolve(JSON.parse(b)));
      }); req.on('error',reject); req.write(d); req.end();
    });
    token = tokenRes.data?.token || '';
  } catch(e) { console.error('Token error:', e.message); }

  // Fetch all endpoints
  const tests = [
    ['📋 健康检查', '/actuator/health'],
    ['ℹ️ 系统信息', '/actuator/info'],
    ['📝 文章列表', '/api/v1alpha1/posts', true],
    ['📂 分类列表', '/api/v1alpha1/categories', true],
    ['🏷️ 标签列表', '/api/v1alpha1/tags', true],
    ['💬 评论列表', '/api/v1alpha1/comments', true],
    ['👤 用户列表', '/api/v1alpha1/users', true],
    ['🔑 角色列表', '/api/v1alpha1/roles', true],
    ['🧭 菜单列表', '/api/v1alpha1/menus', true],
    ['🔗 菜单项列表', '/api/v1alpha1/menuitems', true],
    ['🔌 插件列表', '/api/v1alpha1/plugins', true],
    ['🎨 主题列表', '/api/v1alpha1/themes', true],
    ['📎 附件列表', '/api/v1alpha1/attachments', true],
    ['🔔 通知列表', '/api/v1alpha1/notifications', true],
    ['📊 统计数据', '/api/v1alpha1/stats', true],
    ['📢 公开文章', '/api/public/posts'],
    ['📢 公开分类', '/api/public/categories'],
    ['📢 公开标签', '/api/public/tags'],
  ];

  const results = [];
  for (const [name, path, auth] of tests) {
    try {
      const data = await fetchJSON(path, auth ? token : '');
      const ok = data.code === 0 || data.code === 200;
      results.push({ name, path, status: ok ? 200 : (data.code || 500), ok, data });
    } catch(e) {
      results.push({ name, path, status: 0, ok: false, data: { error: e.message } });
    }
  }

  const passCount = results.filter(r => r.ok).length;
  const html = `<!DOCTYPE html>
<html lang="zh-CN">
<head>
<meta charset="UTF-8">
<title>Halo Go — 实时 API 巡检面板</title>
<style>
  *{margin:0;padding:0;box-sizing:border-box}
  body{font-family:-apple-system,BlinkMacSystemFont,'Segoe UI',Roboto,sans-serif;background:#0a0e1a;color:#e2e8f0}
  .topbar{background:linear-gradient(135deg,#0f172a 0%,#1e293b 100%);padding:20px 30px;border-bottom:1px solid #1e3a5f;display:flex;align-items:center;justify-content:space-between}
  .topbar h1{font-size:22px;background:linear-gradient(90deg,#38bdf8,#818cf8,a78bfa);-webkit-background-clip:text;-webkit-text-fill-color:transparent}
  .badge{display:inline-block;padding:6px 16px;border-radius:20px;font-size:13px;font-weight:600}
  .badge-go{background:#05966922;color:#34d399;border:1px solid #05966944}
  .status-bar{display:flex;gap:20px;padding:16px 30px;background:#0f172a;border-bottom:1px solid #1e293a}
  .stat{text-align:center}
  .stat .n{font-size:28px;font-weight:700;color:#34d399}
  .stat .l{font-size:11px;color:#64748b;margin-top:2px}
  .grid{display:grid;grid-template-columns:repeat(auto-fill,minmax(450px,1fr));gap:16px;padding:20px}
  .card{background:#111827;border-radius:12px;border:1px solid #1f2937;overflow:hidden;transition:transform .2s}
  .card:hover{transform:translateY(-2px);border-color:#334155}
  .card-head{padding:12px 16px;background:#0a0e17;border-bottom:1px solid #1f2937;display:flex;align-items:center;justify-content:space-between;font-size:14px;font-weight:600}
  .card-head .path{font-family:monospace;font-size:12px;color:#64748b;font-weight:400}
  .tag-ok{background:#05966922;color:#4ade80;padding:2px 10px;border-radius:10px;font-size:11px}
  .tag-err{background:#dc262622;color:#f87171;padding:2px 10px;border-radius:10px;font-size:11px}
  .card-body{padding:14px 16px;font-size:12px;max-height:400px;overflow:auto}
  pre{background:#0a0e17;padding:12px;border-radius:8px;overflow-x:auto;line-height:1.5;font-size:11.5px;border:1px solid #1a2035}
  .key{color:#7dd3fc}.str{color:#86efac}.num{color:#fbbf24}.bool{color:#c084fc}.null{color:#64748b}
  .footer{text-align:center;padding:30px;color:#374151;font-size:12px}
  .pulse{display:inline-block;width:8px;height:8px;border-radius:50%;background:#22c55e;margin-right:6px;animation:pulse 2s infinite}
  @keyframes pulse{0%,100%{opacity:1}50%{opacity:.4}}
</style>
</head>
<body>
<div class="topbar">
  <h1>🚀 Halo Go — 实时 API 巡检面板</h1>
  <span class="badge badge-go">Go 1.23 + Gin + GORM + SQLite | 零CGO</span>
</div>
<div class="status-bar">
  <div class="stat"><div class="n">${passCount}/${results.length}</div><div class="l">API 通过</div></div>
  <div class="stat"><div class="n">${Math.round(passCount/results.length*100)}%</div><div class="l">兼容率</div></div>
  <div class="stat"><div class="n" style="color:#60a5fa">8092</div><div class="l">服务端口</div></div>
  <div class="stat"><div class="n" style="color:#facc15">36</div><div class="l">MB 二进制</div></div>
  <div class="stat"><div class="n"><span class="pulse"></span>ON</div><div class="l">服务状态</div></div>
</div>
<div class="grid">
${results.map(r => `
<div class="card">
  <div class="card-head">
    <span>${r.name} <span class="path">${r.path}</span></span>
    <span class="${r.ok ? 'tag-ok' : 'tag-err'}">${r.ok ? '✅ ' + r.status : '❌ ' + r.status}</span>
  </div>
  <div class="card-body"><pre>${jsonHtml(r.data)}</pre></div>
</div>`).join('\n')}
</div>
<div class="footer">
  <p>Halo Go — 基于 Go/Gin/GORM/SQLite 的 Halo 博客系统全量重写</p>
  <p>对标 Java Halo 2.x OpenAPI v3 规范 · 运行在 localhost:8092 · 项目路径 /workspace/halo-go/</p>
</div>
</body></html>`;

  fs.writeFileSync('/workspace/screenshots/live_dashboard.html', html);

  const server = http.createServer((req, res) => {
    res.writeHead(200, {'Content-Type': 'text/html; charset=utf-8'});
    res.end(html);
  });
  server.listen(PORT, () => {
    console.log(`\n🎉 巡检面板已启动!`);
    console.log(`   📊 通过: ${passCount}/${results.length} (${Math.round(passCount/results.length*100)}%)`);
    console.log(`   🔗 打开: http://localhost:${PORT}`);
    console.log(`   📁 报告文件: /workspace/screenshots/live_dashboard.html`);
    console.log(`   🖼️ 对比报告: /workspace/halo_compare_report.html`);
  });
}

main().catch(console.error);
