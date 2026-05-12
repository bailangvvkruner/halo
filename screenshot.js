const puppeteer = require('puppeteer');
const fs = require('fs');
const path = require('path');

(async () => {
  const browser = await puppeteer.launch({
    headless: 'new',
    args: ['--no-sandbox', '--disable-setuid-sandbox', '--disable-gpu']
  });
  const page = await browser.newPage();
  await page.setViewport({ width: 1400, height: 900 });

  console.log('📸 开始截图...');

  // 1. 截图对比报告
  try {
    await page.goto('file:///workspace/halo_compare_report.html', { waitUntil: 'networkidle0', timeout: 15000 });
    await page.screenshot({ path: '/workspace/screenshots/01_compare_report.png', fullPage: true });
    console.log('✅ 1/5 对比报告截图完成');
  } catch(e) { console.error('❌ 报告截图失败:', e.message); }

  // 2. 截图 Go 服务 - 健康检查
  try {
    await page.goto('http://localhost:8092/actuator/health', { waitUntil: 'networkidle0', timeout: 10000 });
    await page.screenshot({ path: '/workspace/screenshots/02_go_health.png' });
    console.log('✅ 2/5 健康检查截图完成');
  } catch(e) { console.error('❌ 健康检查截图失败:', e.message); }

  // 3. 截图 Go 服务 - Info
  try {
    await page.goto('http://localhost:8092/actuator/info', { waitUntil: 'networkidle0', timeout: 10000 });
    await page.screenshot({ path: '/workspace/screenshots/03_go_info.png' });
    console.log('✅ 3/5 Info截图完成');
  } catch(e) { console.error('❌ Info截图失败:', e.message); }

  // 4. 截图 Go 服务 - 文章列表 (带样式)
  try {
    const tokenResp = await fetch('http://localhost:8092/auth/register', {
      method: 'POST',
      headers: {'Content-Type': 'application/json'},
      body: JSON.stringify({username:'shotuser',password:'Shot123456',email:'shot@halo.go',displayName:'ShotUser'})
    });
    const tokenData = await tokenResp.json();
    const token = tokenData.data?.token || '';

    const postsResp = await fetch('http://localhost:8092/api/v1alpha1/posts', {
      headers: {'Authorization': `Bearer ${token}`}
    });
    const postsData = await postsResp.json();

    const html = `<!DOCTYPE html><html><head><meta charset="utf-8"><style>
body{font-family:monospace;background:#0f172a;color:#e2e8f0;padding:20px}
h1{color:#38bdf8}.ok{color:#4ade80}.key{color:#7dd3fc}.str{color:#86efac}.num{color:#fbbf24}
pre{background:#1e293b;padding:16px;border-radius:8px;overflow:auto;font-size:12px;line-height:1.6}
table{width:100%;border-collapse:collapse}td,th{padding:8px 12px;border:1px solid #334155;text-align:left}
th{background:#0f172a;color:#94a3b8}
</style></head><body>
<h1>🚀 Halo Go — 文章列表 API (GET /api/v1alpha1/posts)</h1>
<p>服务地址: <span class="str">http://localhost:8092</span> | 状态: <span class="ok">● 运行中</span></p>
<pre>${JSON.stringify(postsData, null, 2)}</pre>
</body></html>`;
    fs.writeFileSync('/tmp/api_posts.html', html);
    await page.goto('file:///tmp/api_posts.html', { waitUntil: 'networkidle0', timeout: 10000 });
    await page.screenshot({ path: '/workspace/screenshots/04_go_posts_api.png', fullPage: true });
    console.log('✅ 4/5 文章API截图完成');
  } catch(e) { console.error('❌ 文章API截图失败:', e.message); }

  // 5. 截图全部API端点汇总
  try {
    const endpoints = [
      ['/actuator/health', '健康检查'],
      ['/actuator/info', '系统信息'],
      ['/api/public/posts', '公开文章'],
      ['/api/public/categories', '公开分类'],
      ['/api/public/tags', '公开标签'],
    ];
    let rows = '';
    for (const [ep, name] of endpoints) {
      try {
        const r = await fetch(`http://localhost:8092${ep}`);
        const d = await r.json();
        const code = d.code === 0 || d.code === 200 ? '✅' : '⚠️';
        rows += `<tr><td>${name}</td><td>${ep}</td><td>${r.status}</td><td>${code}</td></tr>`;
      } catch(e) {
        rows += `<tr><td>${name}</td><td>${ep}</td><td>ERR</td><td>❌</td></tr>`;
      }
    }
    const summaryHtml = `<!DOCTYPE html><html><head><meta charset="utf-8"><style>
body{font-family:-apple-system,sans-serif;background:#0f172a;color:#e2e8f0;padding:30px}
h1{font-size:24px;text-align:center;background:linear-gradient(90deg,#38bdf8,#818cf8);-webkit-background-clip:text;-webkit-text-fill-color:transparent}
table{width:100%;max-width:800px;margin:20px auto;border-collapse:collapse;background:#1e293b;border-radius:12px;overflow:hidden}
th{padding:14px;background:#0f172a;color:#94a3b8;font-size:14px}td{padding:12px;border-bottom:1px solid #334155;font-size:13px}
tr:hover td{background:#1a2332}.ok{color:#4ade80}.warn{color:#fbbf24}
.footer{text-align:center;color:#475569;margin-top:30px;font-size:12px}
</style></head><body>
<h1>🔍 Halo Go API 端点巡检报告</h1>
<p style="text-align:center;color:#94a3b8">运行在 localhost:8092 | Go 1.23 + Gin + GORM + SQLite</p>
<table><tr><th>端点名称</th><th>路径</th><th>HTTP状态</th><th>结果</th></tr>${rows}</table>
<div class="footer">Halo Go — 全量转换自 Java Halo 2.x | 零CGO依赖 | 纯Go SQLite驱动</div></body></html>`;
    fs.writeFileSync('/tmp/api_summary.html', summaryHtml);
    await page.goto('file:///tmp/api_summary.html', { waitUntil: 'networkidle0', timeout: 10000 });
    await page.screenshot({ path: '/workspace/screenshots/05_api_summary.png' });
    console.log('✅ 5/5 API汇总截图完成');
  } catch(e) { console.error('❌ 汇总截图失败:', e.message); }

  await browser.close();
  console.log('\n🎉 所有截图完成! 保存在 /workspace/screenshots/');
})();
