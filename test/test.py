import urllib.request
import threading
import json
import time

N = 10000 #总请求数
C = 100 # 并发数

headers = {
        'User-Agent':'Mozilla/5.0 (X11; Ubuntu; Linux x86_64; rv:94.0) Gecko/20100101 Firefox/94.0'
    }

lock = threading.Lock()

def crawler(url):
    global N,headers, lock
    data = {
        "uid":123
    }
    jsonData = json.loads(str(data).replace('\'', '\"'))

    while N > 0:
        try:
            req = urllib.request.Request(url=url, headers=headers)
            response = urllib.request.urlopen(req,json.dumps(jsonData).encode('utf-8')).read()
        except Exception as e:
            print(e)
        lock.acquire()
        N -= 1
        lock.release()


def main():
    assert N >= C # 总请求数应大于并发数
    thread = []
    for i in range(C):
        url =  'http://127.0.0.1:8080/api/get_wallet_list'
        t = threading.Thread(target = crawler, args = (url,))
        thread.append(t)
    start = time.time()
    for t in thread:
        t.start()
    for t in thread:
        t.join()
    end = time.time()
    print(f'{end - start}s')


if __name__ == '__main__':
    main()