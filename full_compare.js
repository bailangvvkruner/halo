const { chromium } = require('playwright');
const http = require('http');

const JAVA_BASE = 'http://localhost:8091';
const GO_BASE = 'http://localhost:8090';
const ADMIN = { username: 'admin', password: 'admin', email: 'admin@halo.run', nickname: 'admin' };

function apiPost(url, body, token) {
  return new Promise((resolve, reject) => {
    const u = new URL(url);
    const headers = { 'Content-Type': 'application/json' };
    if (token) headers['Authorization'] = `Bearer ${token}`;
    const req = http.request({ hostname: u.hostname, port: u.port, path: u.pathname + u.search, method: 'POST', headers }, (res) => {
      let data = '';
      res.on('data', c => data += c);
      res.on('end', () => {
        try { resolve({ status: res.statusCode, body: JSON.parse(data) }); }
        catch { resolve({ status: res.statusCode, body: data }); }
      });
    });
    req.on('error', reject);
    req.setTimeout(10000, () => { req.destroy(); reject(new Error('timeout')); });
    req.write(JSON.stringify(body));
    req.end();
  });
}

function apiGet(url, token) {
  return new Promise((resolve, reject) => {
    const u = new URL(url);
    const headers = token ? { 'Authorization': `Bearer ${token}` } : {};
    http.get({ hostname: u.hostname, port: u.port, path: u.pathname + u.search, headers }, (res) => {
      let data = '';
      res.on('data', c => data += c);
      res.on('end', () => {
        try { resolve({ status: res.statusCode, body: JSON.parse(data) }); }
        catch { resolve({ status: res.statusCode, body: data }); }
      });
    }).on('error', reject);
  });
}

async function initJavaHalo(page) {
  console.log('🔧 初始化 Java Halo...');
  
  await page.goto(`${JAVA_BASE}/console/login`, { waitUntil: 'networkidle', timeout: 30000 });
  await page.waitForTimeout(3000);

  const url = page.url();
  console.log(`   当前URL: ${url}`);

  if (url.includes('/setup') || url.includes('/initialize')) {
    console.log('   检测到初始化向导...');
    
    await page.waitForSelector('input', { timeout: 10000 });
    
    const inputs = await page.$$('input');
    for (const input of inputs) {
      const name = await input.getAttribute('name');
      const placeholder = await input.getAttribute('placeholder');
      const type = await input.getAttribute('type');
      
      if (name === 'siteTitle' || placeholder?.includes('站点') || placeholder?.includes('title')) {
        await input.fill('Halo Test Site');
      } else if (name === 'username' || placeholder?.includes('用户') || placeholder?.includes('username')) {
        await input.fill(ADMIN.username);
      } else if (name === 'password' || placeholder?.includes('密码') || placeholder?.includes('password')) {
        await input.fill(ADMIN.password);
      } else if (name === 'email' || placeholder?.includes('邮箱') || placeholder?.includes('email')) {
        await input.fill(ADMIN.email);
      } else if (name === 'nickname' || placeholder?.includes('昵称') || placeholder?.includes('nickname')) {
        await input.fill(ADMIN.nickname);
      }
    }

    await page.screenshot({ path: '/workspace/screenshots/java_setup.png', fullPage: true });
    console.log('   📸 截图: java_setup.png');

    const buttons = await page.$$('button');
    for (const btn of buttons) {
      const text = await btn.textContent();
      if (text.includes('初始化') || text.includes('下一步') || text.includes('完成') || text.includes('提交')) {
        await btn.click();
        console.log(`   点击: ${text.trim()}`);
        await page.waitForTimeout(3000);
        break;
      }
    }

    await page.waitForTimeout(5000);
    console.log(`   初始化后URL: ${page.url()}`);
  }

  return page.url();
}

async function loginJavaConsole(page) {
  console.log('🔑 登录 Java Halo 控制台...');
  
  await page.goto(`${JAVA_BASE}/console/login`, { waitUntil: 'networkidle', timeout: 30000 });
  await page.waitForTimeout(3000);

  const url = page.url();
  console.log(`   当前URL: ${url}`);

  if (url.includes('/login') || url.includes('/console')) {
    await page.waitForSelector('input', { timeout: 10000 }).catch(() => {});
    
    const inputs = await page.$$('input');
    for (const input of inputs) {
      const placeholder = await input.getAttribute('placeholder');
      const type = await input.getAttribute('type');
      
      if (placeholder?.includes('用户名') || placeholder?.includes('账号') || placeholder?.includes('username')) {
        await input.fill(ADMIN.username);
      } else if (type === 'password' || placeholder?.includes('密码')) {
        await input.fill(ADMIN.password);
      }
    }

    const buttons = await page.$$('button');
    for (const btn of buttons) {
      const text = await btn.textContent();
      if (text.includes('登录') || text.includes('Login') || text.includes('登 录')) {
        await btn.click();
        console.log(`   点击登录...`);
        await page.waitForTimeout(5000);
        break;
      }
    }
  }

  console.log(`   登录后URL: ${page.url()}`);
  return page.url();
}

async function compareWithAPI() {
  console.log('\n═══════════════════════════════════════════════');
  console.log('  API 数据对比 (带认证)');
  console.log('═══════════════════════════════════════════════\n');

  const goToken = (await apiPost(`${GO_BASE}/auth/login`, ADMIN)).body?.data?.token;
  console.log(`Go Token: ${goToken ? goToken.substring(0, 25) + '...' : 'FAILED'}`);

  if (!goToken) {
    const regRes = await apiPost(`${GO_BASE}/auth/register`, { ...ADMIN, email: ADMIN.email });
    const newToken = regRes.body?.data?.token;
    console.log(`Go Register Token: ${newToken ? newToken.substring(0, 25) + '...' : 'FAILED'}`);
  }

  const endpoints = [
    { path: '/api/v1alpha1/posts?page=1&size=10', name: '文章列表' },
    { path: '/api/v1alpha1/categories?page=1&size=10', name: '分类列表' },
    { path: '/api/v1alpha1/tags?page=1&size=10', name: '标签列表' },
    { path: '/api/v1alpha1/menus?page=1&size=10', name: '菜单列表' },
    { path: '/api/v1alpha1/menuitems?page=1&size=10', name: '菜单项列表' },
    { path: '/api/v1alpha1/singlepages?page=1&size=10', name: '单页列表' },
    { path: '/api/v1alpha1/stats', name: '统计数据' },
  ];

  let pass = 0, fail = 0;

  for (const ep of endpoints) {
    const [goRes, javaRes] = await Promise.all([
      apiGet(`${GO_BASE}${ep.path}`, goToken),
      apiGet(`${JAVA_BASE}${ep.path}`).catch(() => ({ status: 0, body: {} }))
    ]);

    const goOk = goRes.status === 200;
    const goCount = goRes.body?.data?.items?.length || goRes.body?.data?.total || goRes.body?.items?.length || '?';
    
    console.log(`   ${goOk ? '✅' : '❌'} ${ep.name}: Go=${goRes.status}(${goCount}条) Java=${javaRes.status}`);
    
    if (goOk) pass++; else fail++;
  }

  console.log(`\n   📊 Go API: ${pass}/${pass+fail} 端点正常`);
}

async function main() {
  console.log('═══════════════════════════════════════════════');
  console.log('  🚀 Halo Go vs Java — 全面对比测试');
  console.log('═══════════════════════════════════════════════\n');

  const browser = await chromium.launch({ headless: true });

  try {
    const ctx = await browser.newContext({ viewport: { width: 1440, height: 900 } });
    const page = await ctx.newPage();

    await initJavaHalo(page);
    
    await page.screenshot({ path: '/workspace/screenshots/java_after_init.png', fullPage: true });
    console.log('   📸 截图: java_after_init.png');

    // Navigate to console pages
    const pages_to_visit = [
      { path: '/console/dashboard', name: '仪表盘' },
      { path: '/console/posts', name: '文章管理' },
      { path: '/console/single-pages', name: '页面管理' },
      { path: '/console/themes', name: '主题管理' },
      { path: '/console/menus', name: '菜单管理' },
    ];

    for (const p of pages_to_visit) {
      try {
        await page.goto(`${JAVA_BASE}${p.path}`, { waitUntil: 'networkidle', timeout: 20000 });
        await page.waitForTimeout(2000);
        await page.screenshot({ path: `/workspace/screenshots/java_${p.name}.png`, fullPage: true });
        console.log(`   📸 ${p.name}: java_${p.name}.png`);
      } catch (e) {
        console.log(`   ⚠️ ${p.name}: ${e.message.substring(0, 80)}`);
      }
    }

    // Try post editor
    try {
      await page.goto(`${JAVA_BASE}/console/posts/editor`, { waitUntil: 'networkidle', timeout: 20000 });
      await page.waitForTimeout(2000);
      await page.screenshot({ path: '/workspace/screenshots/java_post_editor.png', fullPage: true });
      console.log('   📸 文章编辑器: java_post_editor.png');
    } catch (e) {
      console.log(`   ⚠️ 文章编辑器: ${e.message.substring(0, 80)}`);
    }

    await ctx.close();

    // API comparison
    await compareWithAPI();

  } catch (e) {
    console.error('测试异常:', e.message);
  } finally {
    await browser.close();
  }

  console.log('\n═══════════════════════════════════════════════');
  console.log('  ✅ 全面对比测试完成');
  console.log('  截图保存在 /workspace/screenshots/');
  console.log('═══════════════════════════════════════════════');
}

main();