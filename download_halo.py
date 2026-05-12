#!/usr/bin/env python3
import requests, os, sys, time

JAR_PATH = '/workspace/halo-2.19.0.jar'
URL = 'https://github.com/halo-dev/halo/releases/download/v2.19.0/halo-2.19.0.jar'
HEADERS = {'User-Agent': 'Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 Chrome/120.0.0.0 Safari/537.36'}

total_size = None
start_time = time.time()

def get_total_size():
    r = requests.head(URL, headers=HEADERS, allow_redirects=True, timeout=15)
    return int(r.headers.get('content-length', 0))

def download_chunk(offset=0):
    headers = dict(HEADERS)
    if offset > 0:
        headers['Range'] = f'bytes={offset}-'
    r = requests.get(URL, headers=headers, stream=True, timeout=(30, 300), allow_redirects=True)
    return r

print("获取文件大小...")
total_size = get_total_size()
print(f"文件大小: {total_size/1024/1024:.1f} MB")

existing = os.path.getsize(JAR_PATH) if os.path.exists(JAR_PATH) else 0
print(f"已有: {existing/1024/1024:.1f} MB")

mode = 'ab' if existing > 0 else 'wb'
downloaded = existing
attempts = 0
MAX_RETRIES = 50

while downloaded < total_size and attempts < MAX_RETRIES:
    attempts += 1
    try:
        r = download_chunk(downloaded)
        print(f"\n连接建立，偏移={downloaded}, 状态={r.status_code}")
        
        with open(JAR_PATH, mode) as f:
            for chunk in r.iter_content(chunk_size=128*1024):
                if chunk:
                    f.write(chunk)
                    downloaded += len(chunk)
                    pct = downloaded*100/total_size
                    elapsed = time.time() - start_time
                    speed = downloaded/elapsed if elapsed > 0 else 0
                    bar_len = 30
                    filled = int(bar_len * pct / 100)
                    bar = '█' * filled + '░' * (bar_len - filled)
                    sys.stdout.write(f'\r  [{bar}] {pct:.1f}%  {downloaded/1024/1024:.1f}/{total_size/1024/1024:.1f}MB  {speed/1024:.0f}KB/s  尝试{attempts}',)
                    sys.stdout.flush()
        mode = 'ab'
        print(f"\n完成！ {downloaded/1024/1024:.1f}MB")
        break
    except Exception as e:
        print(f"\n错误 #{attempts}: {e}")
        if downloaded > 0:
            mode = 'ab'
        time.sleep(3)

print(f"\n最终: {downloaded/1024/1024:.1f}MB / {total_size/1024/1024:.1f}MB")
if downloaded >= total_size * 0.99:
    print("✅ 下载完成！")
else:
    print(f"⚠️ 未完成 ({downloaded/total_size*100:.1f}%)")
