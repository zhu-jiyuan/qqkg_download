import time
import asyncio
import httpx
from tqdm.auto import tqdm
import sys
import re
import os

path = r'./music/'
if not os.path.exists(path):
    os.makedirs(path)


async def download_music(client: httpx.AsyncClient, shareid: str, title: str):
    data_url = "http://cgi.kg.qq.com/fcgi-bin/fcg_get_play_url?shareid=" + str(shareid)

    file = path + title + '.m4a'

    async with client.stream("GET", data_url, timeout = 5) as response:
        with open(file, 'wb') as f:
            with tqdm(
                unit='B', unit_scale=True, unit_divisor=1024, miniters=1,
                desc=title, total=int(response.headers.get('content-length', 0))
            ) as pbar:
                async for chunk in response.aiter_bytes(chunk_size=4096):
                    f.write(chunk)
                    pbar.update(len(chunk))
        

async def get_music_url_list(client: httpx.AsyncClient, uid: str):
    start_count = 1
    ugc_list_all = []

    headers = {
        "user-agent": "Mozilla/5.0 (iPhone; CPU iPhone OS 16_6 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/16.6 Mobile/15E148 Safari/604.1",
        "referer": "https://static-play.kg.qq.com/node/personal_v2/?uid=" + str(uid),
    }

    async def _get_url_list(start_count):
        ts_ms = int(time.time() * 1000)
        url = "https://node.kg.qq.com/fcgi-bin/kg_ugc_get_homepage?outCharset=utf-8&from=1&nocache={}&format=json&type=get_uinfo&start={}&num=10&share_uid={uid}&g_tk=1164660242&g_tk_openkey=1164660242".format(ts_ms, start_count, uid=uid)

        resp = await client.get(url, headers=headers)

        resp = resp.json()
        ugc_list = resp["data"]["ugclist"]
        return ugc_list


    while True:
        ugc_list = await _get_url_list(start_count)
        ugc_list_all.extend(ugc_list)
        if len(ugc_list) < 10:
            break
        start_count += 1


    return ugc_list_all


async def main(uid):
     async with httpx.AsyncClient(follow_redirects=True) as client:
        ugclist = await get_music_url_list(client, uid)
        tasks = []
        for i in range(len(ugclist)):
            ugc = ugclist[i]
            title = ugc["title"]
            shareid = ugc["shareid"]
            tasks.append(download_music(client, shareid, title))
            if i%10 == 0:
                await asyncio.gather(*tasks)
                tasks = []

        if len(tasks) > 0:
            await asyncio.gather(*tasks)



def get_uid(url):
    uid_match = re.search(r'uid=([0-9a-f]+)', url)
    if uid_match==None:
        print('uid not found')
        print("Please provide a valid url.\ne.g. https://node.kg.qq.com/personal?uid=1223342")
        sys.exit(1)

    uid = uid_match.group(1)
    return uid

    
if __name__ == "__main__":
    url = sys.argv[1]
    uid = get_uid(url)
    asyncio.run(main(uid))
