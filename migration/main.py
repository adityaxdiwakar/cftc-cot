import requests
import time
import os

types = ["disaggregated", "financials"]
headers = {"Authorization":  os.getenv("API_TOKEN")}
for t in types:
    r = requests.get(f"https://api.adi.wtf/cots/{t}").json()
    for key in r:
        pair = (key, r[key]["name"])
        requests.post(f"https://api.aditya.diwakar.io/cot/{t}", json={
            "id": pair[0],
            "name": pair[1]
        }, headers=headers)

        for date in r[key]["data"]:
            n_date = int(date.replace("-", ""))
            d = r[key]["data"][date]
            for field in d:
                d[field] = int(d[field])
            d["date"] = n_date
            requests.put(f"https://api.aditya.diwakar.io/cot/{t}/{key}",
                  json=d,
                  headers=headers)
            time.sleep(0.05)
