from cryptography.fernet import Fernet
import requests
from requests.utils import requote_uri

cyp = "gAAAAABkSfi9pc6mfUlosvRlvVPOpHRU-PszCISzPwxgGNyj2tO_qiUbbGKdqSyWndfxUCHARc4bpgmiObilNXorlR_Wcxo9Mw_wLhD4jSoxua_rzznMbhlzbJkHnichrv6lEyHcZ_KSoVob53XKKz4Y_6lrp-hdqLyIzM0JLf6b3XV2iY6sCTUiZrK34hU9cb-aZZBET_y3PKBlegwOojdBJ-z2lcMMpQmAUGL6Qq7Cqg1teUEy7pZZU0MkGfqfm0PKrOXD_z7YAkR7NDh2w7jvRc4PUy1_HRBxFSA77wojPxIqijUukN4AxrmeQl4enEfHL4ymgjxW6xXlPAla1z3ryXg-eCL0oDaeLhEaxbt4I28_lrFrmM6NVeJs503Rxprw_gt7Bwwfe8IJ08U4DOTvDufvioJ6zcFi27cRDRGPd5MbO4qYg43Myug5yeAJtWlLvvSJnLhO082is-kuDDrQ6sLfmNs3MhMb-ehgfFrORb8t-Ibhtn5i4VJ1OlI3dYZWTNTD6yHuMThwrnksVcIFIgTdSwCWC7bEn0LtWIdQZjcMLaD3ZiHyiVk5j2_8NQg3YEePyLdXHStEJaSk5vhVRoV85DS_cFTgkArrzhwF4FF7wLHLmlYTVLPi04-ijLRKmWOoBE-K-JcbFnI1B-VxnZ-lKHfhk_k3P9eiqa3kuA421AiVAz8KkhweYXdDVIHJA_2FWICwVQdeJ1FHUkm3QnfLJ2cq1vFyvqu84vJ13hkbHLjc1IlueJILlria1uB-ggVy9_HTYE0PVBl3d5p5L1qnkQG0Mi46bCfAHhuDaiil3rlP1HGG9nE5bHJ7TCW_kTuh3fFnLE_9PmjMBOSZD-Ej6-kncxbkinm4uqZYuvBp-HvrmuWng4TQqeZ_JTsO4bDoaDUYYdID7BL29dKTPJDO3SDWPw=="
key = "0IAitdUbAqCG4FXEQDvouB1iyddatp3WZ2igXxCSluw="
f = Fernet(key)
mes = f.decrypt(cyp).decode()
ans = []
for word in mes.split(' '):
    if word.startswith('U+'):
        ans.append(chr(int(word[2::], 16)).encode('utf-8').decode())
    else:
        ans.append(word)

ans = ' '.join(ans)
ans_enc = requote_uri(ans)

requests.post(
    f'https://quest.squadcast.tech/api/310620104093/submit/emoji?answer={ans_enc}',
    open('emoji.py', 'r').read(),
    headers={"Content-Type": 'text/plain'}
)
