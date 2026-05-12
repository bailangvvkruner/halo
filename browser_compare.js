const { chromium } = require('playwright');
const fs = require('fs');
const http = require('http');
const path = require('path');

const SCREENSHOTS_DIR = '/workspace/screenshots';
const JAVA_HALO_URL = 'http://localhost:8091';
const GO_HALO_URL = 'http://localhost:8090';

// 模拟真人操作步骤：每个步骤包含 URL 和操作描述
const TEST_SCENARIOS = [
  {
    name: '🏠 首页（前台）',
    url: '/',
    action: async (page) => {
      await page.waitForLoadState('networkidle');
      await page.waitForTimeout(1500);
    }
  },
  {
    name: '📋 健康检查',
    url: '/actuator/health',
    action: async (page) => {
      await page.waitForLoadState('networkidle');
    }
  },
  {
    name: 'ℹ️ 系统信息',
    url: '/actuator/info',
    action: async (page) => {
      await page.waitForLoadState('networkidle');
    }
  },
  {
    name: '🔐 控制台入口（后台登录）',
    url: '/console',
    action: async (page) => {
      await page.waitForLoadState('networkidle');
      await page.waitForTimeout(2000);
    }
  },
  {
    name: '📝 文章列表 API',
    url: '/api/v1alpha1/posts',
    action: async (page) => {
      await page.waitForLoadState('networkidle');
    }
  },
  {
    name: '📂 分类列表 API',
    url: '/api/v1alpha1/categories',
    action: async (page) => {
      await page.waitForLoadState('networkidle');
    }
  },
  {
    name: '🏷️ 标签列表 API',
    url: '/api/v1alpha1/tags',
    action: async (page) => {
      await page.waitForLoadState('networkidle');
    }
  },
  {
    name: '💬 评论列表 API',
    url: '/api/v1alpha1/comments',
    action: async (page) => {
      await page.waitForLoadState('networkidle');
    }
  },
  {
    name: '👤 用户列表 API',
    url: '/api/v1alpha1/users',
    action: async (page) => {
      await page.waitForLoadState('networkidle');
    }
  },
  {
    name: '📄 单页 API',
    url: '/api/v1alpha1/singlepages',
    action: async (page) => {
      await page.waitForLoadState('networkidle');
    }
  },
  {
    name: '🍽️ 菜单 API',
    url: '/api/v1alpha1/menus',
    action: async (page) => {
      await page.waitForLoadState('networkidle');
    }
  },
  {
    name: '📌 菜单项 API',
    url: '/api/v1alpha1/menuitems',
    action: async (page) => {
      await page.waitForLoadState('networkidle');
    }
  },
  {
    name: '📊 统计信息 API',
    url: '/api/v1alpha1/stats',
    action: async (page) => {
      await page.waitForLoadState('networkidle');
    }
  },
  {
    name: '📎 附件 API',
    url: '/api/v1alpha1/attachments',
    action: async (page) => {
      await page.waitForLoadState('networkidle');
    }
  },
  {
    name: '🔌 插件 API',
    url: '/api/v1alpha1/plugins',
    action: async (page) => {
      await page.waitForLoadState('networkidle');
    }
  },
  {
    name: '🎨 主题 API',
    url: '/api/v1alpha1/themes',
    action: async (page) => {
      await page.waitForLoadState('networkidle');
    }
  },
  {
    name: '👥 角色 API',
    url: '/api/v1alpha1/roles',
    action: async (page) => {
      await page.waitForLoadState('networkidle');
    }
  },
  {
    name: '📸 快照 API',
    url: '/api/v1alpha1/snapshots',
    action: async (page) => {
      await page.waitForLoadState('networkidle');
    }
  },
];

async function checkServiceAlive(url) {
  return new Promise((resolve) => {
    const req = http.get(url, { timeout: 5000 }, (res) => {
      resolve(res.statusCode < 500);
    });
    req.on('error', () => resolve(false));
    req.setTimeout(5000, () => { req.destroy(); resolve(false); });
  });
}

async function takeScreenshot(page, scenario, label) {
  const safeName = scenario.name.replace(/[\/\?<>\\:\*\|"]/g, '_');
  const filename = `${label}_${safeName}.png`;
  const filepath = path.join(SCREENSHOTS_DIR, filename);
  await page.screenshot({ path: filepath, fullPage: true });
  return { filename, filepath };
}

function generateReport(results) {
  const rows = results.map(r => {
    const javaImg = r.java ? `<img src="${r.java.filename}" style="max-width:100%;border-radius:8px;border:2px solid #334155" />` : '<div style="padding:40px;color:#f87171;text-align:center">❌ 服务未启动</div>';
    const goImg = r.go ? `<img src="${r.go.filename}" style="max-width:100%;border-radius:8px;border:2px solid #334155" />` : '<div style="padding:40px;color:#f87171;text-align:center">❌ 服务未启动</div>';
    return `
    <tr>
      <td style="font-weight:600;font-size:14px;vertical-align:top;padding-top:16px">${r.name}</td>
      <td style="background:#0f172a">${javaImg}</td>
      <td style="background:#0f172a">${goImg}</td>
    </tr>`;
  }).join('');

  return `<!DOCTYPE html>
<html lang="zh-CN">
<head>
<meta charset="UTF-8">
<title>Halo Go vs Java Halo — 真人操作逐页对比</title>
<style>
  * { margin:0; padding:0; box-sizing:border-box; }
  body { font-family:-apple-system,BlinkMacSystemFont,'Segoe UI',Roboto,sans-serif; background:#0a0e1a; color:#e2e8f0; padding:20px; }
  .header { text-align:center; padding:30px 20px; background:linear-gradient(135deg,#1e293b,#334155); border-radius:16px; margin-bottom:24px; border:1px solid #475569; }
  .header h1 { font-size:26px; background:linear-gradient(90deg,#38bdf8,#818cf8); -webkit-background-clip:text;-webkit-text-fill-color:transparent; }
  .header .sub { color:#94a3b8; margin-top:8px; font-size:13px; line-height:1.6; }
  .badge { display:inline-block; padding:4px 12px; border-radius:20px; font-size:12px; font-weight:600; margin:3px; }
  .badge-java { background:#2563eb22; color:#60a5fa; border:1px solid #2563eb44; }
  .badge-go { background:#05966922; color:#34d399; border:1px solid #05966944; }
  table { width:100%; border-collapse:separate; border-spacing:0; }
  th { padding:14px 16px; background:#0f172a; color:#94a3b8; font-weight:600; font-size:13px; text-align:left; border-bottom:2px solid #38bdf8; }
  th:first-child { width:180px; }
  td { padding:12px; border-bottom:1px solid #1f2937; vertical-align:top; }
  tr:hover td { background:#111827; }
  img:hover { transform:scale(1.02); transition:transform .3s; cursor:zoom-in; }
  .footer { text-align:center; padding:30px; color:#475569; font-size:12px; line-height:1.8; }
  .note { background:#1e293b; border-left:4px solid #f59e0b; padding:16px 20px; border-radius:8px; margin:20px 0; font-size:13px; color:#fbbf24; }
</style>
</head>
<body>

<div class="header">
  <h1>🚀 Halo Go vs Java Halo — 浏览器真人操作逐页对比</h1>
  <div class="sub">
    <p>使用 Playwright Chromium 真实浏览器渲染 · 模拟真人访问每个页面</p>
    <p>
      <span class="badge badge-java">Java Halo 2.19.0 → localhost:8090</span>
      <span class="badge badge-go">Go Halo 0.1.0 → localhost:8091</span>
    </p>
  </div>
</div>

<div class="note">
  ⚠️ <strong>说明：</strong>此报告通过真实 Chromium 浏览器访问两个服务，对每个页面进行完整渲染后截图。
  对比的是<strong>页面实际渲染效果</strong>而非 JSON 返回值。Go 版本为纯 API 后台，无前端控制台 UI；
  Java 版本包含完整的 Vue.js 前端控制台。
</div>

<table>
  <tr>
    <th>页面 / 操作</th>
    <th>🟦 Java Halo (原版)</th>
    <th>🟩 Go Halo (重写)</th>
  </tr>
  ${rows}
</table>

<div class="footer">
  <p>生成时间: ${new Date().toLocaleString('zh-CN', { timeZone: 'Asia/Shanghai' })}</p>
  <p>测试工具: Playwright (Chromium) | 截图目录: ${SCREENSHOTS_DIR}</p>
  <p>Go 技术栈: Go 1.25 + Gin + GORM + SQLite (零CGO) | Java 技术栈: Java 25 + Spring Boot 3.3.3</p>
</div>

</body></html>`;
}

(async () => {
  console.log('═══════════════════════════════════════════════');
  console.log('  🌐 Halo Go vs Java Halo — 真人操作逐页对比');
  console.log('═══════════════════════════════════════════════\n');

  if (!fs.existsSync(SCREENSHOTS_DIR)) fs.mkdirSync(SCREENSHOTS_DIR, { recursive: true });

  // 检查服务状态
  console.log('⏳ 检查服务状态...');
  const javaAlive = await checkServiceAlive(`${JAVA_HALO_URL}/actuator/health`);
  const goAlive = await checkServiceAlive(`${GO_HALO_URL}/actuator/health`);
  console.log(`   🟦 Java Halo (8090): ${javaAlive ? '✅ 运行中' : '❌ 未启动'}`);
  console.log(`   🟩 Go Halo  (8091): ${goAlive ? '✅ 运行中' : '❌ 未启动'}`);
  console.log('');

  let browser;
  try {
    browser = await chromium.launch({
      headless: true,
      args: ['--no-sandbox', '--disable-setuid-sandbox', '--disable-gpu', '--font-render-hinting=none']
    });
  } catch(e) {
    console.error('❌ Playwright Chromium 启动失败:', e.message);
    process.exit(1);
  }

  const results = [];

  for (let i = 0; i < TEST_SCENARIOS.length; i++) {
    const scenario = TEST_SCENARIOS[i];
    console.log(`\n[${i + 1}/${TEST_SCENARIOS.length}] ${scenario.name}`);

    const result = { name: scenario.name, java: null, go: null };

    // 测试 Java Halo
    if (javaAlive) {
      const javaPage = await browser.newPage({ viewport: { width: 1400, height: 900 } });
      try {
        await javaPage.goto(`${JAVA_HALO_URL}${scenario.url}`, { waitUntil: 'networkidle', timeout: 30000 });
        await scenario.action(javaPage);
        result.java = await takeScreenshot(javaPage, scenario, 'java');
        console.log(`   🟦 Java ✅ ${result.java.filename}`);
      } catch(e) {
        console.log(`   🟦 Java ⚠️ ${e.message.slice(0, 60)}`);
      } finally {
        await javaPage.close();
      }
    }

    // 测试 Go Halo
    if (goAlive) {
      const goPage = await browser.newPage({ viewport: { width: 1400, height: 900 } });
      try {
        await goPage.goto(`${GO_HALO_URL}${scenario.url}`, { waitUntil: 'networkidle', timeout: 30000 });
        await scenario.action(goPage);
        result.go = await takeScreenshot(goPage, scenario, 'go');
        console.log(`   🟩 Go  ✅ ${result.go.filename}`);
      } catch(e) {
        console.log(`   🟩 Go  ⚠️ ${e.message.slice(0, 60)}`);
      } finally {
        await goPage.close();
      }
    }

    results.push(result);
  }

  await browser.close();

  // 生成 HTML 报告
  const html = generateReport(results);
  const reportPath = '/workspace/screenshots/browser_compare_report.html';
  fs.writeFileSync(reportPath, html);

  console.log('\n═══════════════════════════════════════════════');
  console.log('  ✅ 全部截图完成！');
  console.log(`  📊 报告: ${reportPath}`);
  console.log(`  📁 截图: ${SCREENSHOTS_DIR}/`);
  console.log(`  🟦 Java 页面: ${results.filter(r => r.java).length} 个`);
  console.log(`  🟩 Go  页面: ${results.filter(r => r.go).length} 个`);
  console.log('═══════════════════════════════════════════════');

  // 启动 HTTP 服务查看报告
  const server = http.createServer((req, res) => {
    if (req.url === '/' || req.url === '/report') {
      res.writeHead(200, { 'Content-Type': 'text/html; charset=utf-8' });
      res.end(html);
    } else {
      res.writeHead(404);
      res.end('Not found');
    }
  });
  server.listen(8899, () => {
    console.log(`\n🌐 对比报告已启动: http://localhost:8899/report`);
  });
})();
