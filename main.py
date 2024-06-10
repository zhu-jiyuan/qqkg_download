import time
import asyncio
import httpx
import sys
import re
import os
import tqdm.asyncio
# from pprint import pprint

path = r'./music/'
if not os.path.exists(path):
    os.makedirs(path)


async def download_music(client: httpx.AsyncClient, shareid: str, title: str, nickname:str):
    data_url = "http://cgi.kg.qq.com/fcgi-bin/fcg_get_play_url?shareid=" + str(shareid)

    file = path + title + '-' + str(nickname) + '.m4a'
    if os.path.exists(file):
        print(f"File {file} already exists")
        return

    async with client.stream("GET", data_url, timeout = 5) as response:
        with open(file, 'wb') as f:
            with tqdm.asyncio.tqdm(
                unit='B', unit_scale=True, unit_divisor=1024, miniters=1,
                desc=title + '-' + nickname, total=int(response.headers.get('content-length', 0))
            ) as pbar:
                async for chunk in response.aiter_bytes(chunk_size=4096):
                    f.write(chunk)
                    pbar.update(len(chunk))
        

async def get_music_data(client: httpx.AsyncClient, uid: str):
    start_count = 1
    ugc_list_all = []

    headers = {
        "user-agent": "Mozilla/5.0 (iPhone; CPU iPhone OS 16_6 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/16.6 Mobile/15E148 Safari/604.1",
        "referer": "https://static-play.kg.qq.com/node/personal_v2/?uid=" + str(uid),
    }

    ts_ms = int(time.time() * 1000)
    url = "https://node.kg.qq.com/fcgi-bin/kg_ugc_get_homepage?outCharset=utf-8&from=1&nocache={}&format=json&type=get_uinfo&start={}&num=10&share_uid={uid}&g_tk=1164660242&g_tk_openkey=1164660242".format(ts_ms, start_count, uid=uid)

    resp = await client.get(url, headers=headers)
    data = resp.json()["data"]
    ugc_total_count = data["ugc_total_count"]
    nickname = data["nickname"]

    async def _get_url_list(start_count):
        ts_ms = int(time.time() * 1000)
        url = "https://node.kg.qq.com/fcgi-bin/kg_ugc_get_homepage?outCharset=utf-8&from=1&nocache={}&format=json&type=get_uinfo&start={}&num=10&share_uid={uid}&g_tk=1164660242&g_tk_openkey=1164660242".format(ts_ms, start_count, uid=uid)

        resp = await client.get(url, headers=headers)

        resp = resp.json()
        ugc_list = resp["data"]["ugclist"]
        ugc_list_all.extend(ugc_list)
        return ugc_list


    get_count = ugc_total_count // 10
    if ugc_total_count % 10 !=0 :
        get_count = get_count + 1

    tasks = []
    for i in range(1, get_count + 1):
        tasks.append(_get_url_list(i))

    await asyncio.gather(*tasks)

    return ugc_list_all, nickname


async def main(uid):
     async with httpx.AsyncClient(follow_redirects=True) as client:
        ugclist, nickname = await get_music_data(client, uid)
        print(f"Found {len(ugclist)} songs by {nickname}")
        tasks = []
        for i in range(len(ugclist)):
            ugc = ugclist[i]
            title = ugc["title"]
            shareid = ugc["shareid"]
            tasks.append(download_music(client, shareid, title, nickname))
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
