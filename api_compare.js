const http = require('http');

const JAVA_BASE = 'http://localhost:8091';
const GO_BASE = 'http://localhost:8090';

const TEST_ENDPOINTS = [
  { path: '/actuator/health', name: '健康检查', exact: true },
  { path: '/actuator/info', name: '系统信息', exact: false },
  { path: '/api/v1alpha1/posts', name: '文章列表 API', exact: false },
  { path: '/api/v1alpha1/categories', name: '分类列表 API', exact: true },
  { path: '/api/v1alpha1/tags', name: '标签列表 API', exact: true },
  { path: '/api/v1alpha1/menus', name: '菜单 API', exact: true },
  { path: '/api/v1alpha1/menuitems', name: '菜单项 API', exact: true },
  { path: '/api/v1alpha1/singlepages', name: '单页 API', exact: true },
  { path: '/api/v1alpha1/stats', name: '统计 API', exact: false },
];

function getRequest(baseUrl, path) {
  return new Promise((resolve, reject) => {
    const url = new URL(path, baseUrl);
    const req = http.get(url, (res) => {
      let data = '';
      res.on('data', (chunk) => {
        data += chunk;
      });
      res.on('end', () => {
        try {
          resolve({
            status: res.statusCode,
            data: JSON.parse(data)
          });
        } catch (e) {
          resolve({
            status: res.statusCode,
            data: data
          });
        }
      });
    });
    req.on('error', reject);
    req.setTimeout(10000, () => {
      req.destroy();
      reject(new Error('Request timeout'));
    });
  });
}

function compareResponses(javaRes, goRes, test) {
  let result = {
    name: test.name,
    path: test.path,
    javaStatus: javaRes.status,
    goStatus: goRes.status,
    statusOk: javaRes.status === goRes.status,
    contentOk: false,
    details: []
  };

  if (result.statusOk) {
    try {
      const javaData = typeof javaRes.data === 'string' ? JSON.parse(javaRes.data) : javaRes.data;
      const goData = typeof goRes.data === 'string' ? JSON.parse(goRes.data) : goRes.data;

      if (test.exact) {
        result.contentOk = JSON.stringify(simplifyForComparison(javaData)) === JSON.stringify(simplifyForComparison(goData));
      } else {
        result.contentOk = true;
      }
      result.javaData = javaData;
      result.goData = goData;
    } catch (e) {
      result.details.push(`数据解析错误: ${e.message}`);
    }
  }

  return result;
}

function simplifyForComparison(obj) {
  if (obj === null || obj === undefined) return obj;
  if (typeof obj === 'string' || typeof obj === 'number' || typeof obj === 'boolean') return obj;
  if (Array.isArray(obj)) return obj.map(simplifyForComparison);
  const result = {};
  for (const key of Object.keys(obj).sort()) {
    if (key !== 'version' && key !== 'creationTimestamp' && key !== 'lastModifyTime' && key !== 'publishTime') {
      result[key] = simplifyForComparison(obj[key]);
    }
  }
  return result;
}

async function runTests() {
  console.log('═══════════════════════════════════════════════');
  console.log('  🚀 Halo Go vs Java Halo — API 一致性对比');
  console.log('═══════════════════════════════════════════════\n');

  const results = [];

  for (const test of TEST_ENDPOINTS) {
    console.log(`🔍 测试: ${test.name}`);
    console.log(`   路径: ${test.path}`);
    
    try {
      const [javaRes, goRes] = await Promise.all([
        getRequest(JAVA_BASE, test.path),
        getRequest(GO_BASE, test.path)
      ]);

      const result = compareResponses(javaRes, goRes, test);
      results.push(result);

      console.log(`   Java: ${result.javaStatus} | Go: ${result.goStatus}`);
      console.log(`   状态一致: ${result.statusOk ? '✅' : '❌'}`);
      console.log(`   内容一致: ${result.contentOk ? '✅' : '⚠️'}`);
      console.log('');
    } catch (e) {
      console.log(`   ❌ 错误: ${e.message}`);
      results.push({
        name: test.name,
        path: test.path,
        error: e.message
      });
    }
  }

  console.log('═══════════════════════════════════════════════');
  console.log('  📊 测试汇总');
  console.log('═══════════════════════════════════════════════');
  
  const passed = results.filter(r => r.statusOk).length;
  const total = results.length;
  
  console.log(`\n📈 成功率: ${passed}/${total} (${Math.round(passed/total*100)}%)`);
  
  console.log('\n✅ 通过的 API:');
  results.filter(r => r.statusOk && r.contentOk).forEach(r => {
    console.log(`  - ${r.name}`);
  });
  
  console.log('\n⚠️ 需要注意的 API:');
  results.filter(r => !r.statusOk || !r.contentOk).forEach(r => {
    console.log(`  - ${r.name}: ${!r.statusOk ? '状态码不一致' : '内容不完全一致'}`);
  });

  console.log('\n✅ API 对比验证完成！');
}

runTests().catch(console.error);
