from dotenv import load_dotenv
load_dotenv()

import threading
import requests
import time
import os

# log start time
start = time.time()

# request for the comma delim'd file from CFTC
disagg_req = requests.get("https://www.cftc.gov/dea/newcot/f_disagg.txt")
fin_req = requests.get("https://www.cftc.gov/dea/newcot/FinFutWk.txt")

# checking if either give non-200 status codes
if disagg_req.status_code != 200 or fin_req.status_code != 200:
    raise Exception("Non-200 Status Code")

# combining lines from both of the requests
disagg_lines = disagg_req.text.split("\r\n")
fin_lines = fin_req.text.split("\r\n")
lines = disagg_lines + fin_lines

# iterate over the lines from the combined list
for index, line in enumerate(lines):

    # extract the data from the lines
    try:
        trimmed_line = line[1:]
        first_quote = trimmed_line.find('"')
        trimmed_line = line[first_quote + 3:]
        splits = trimmed_line.split(",")

        cot_id, name, date = str(splits[2]), line[1:first_quote+1], splits[1]
    except:
        # in the case we are at a header line or something else goes wrong
        continue

    record = {}
    if index < len(disagg_lines):
        record = {
            "prod_long": splits[7],
            "prod_short": splits[8],
            "swap_long": splits[9],
            "swap_short": splits[10],
            "mm_long": splits[12],
            "mm_short": splits[13],
            "other_long": splits[15],
            "other_short": splits[16]
        }
    else:
        record = {
            "deal_long": splits[7],
            "deal_short": splits[8],
            "asset_long": splits[10],
            "asset_short": splits[11],
            "lev_long": splits[13],
            "lev_short": splits[14],
            "other_long": splits[16],
            "other_short": splits[17]
        }

    cot_class = ["disaggregated", "financials"][index >= len(disagg_lines)]

    for key in record:
        record[key] = int(record[key].replace(' ', ''))
    record["date"] = int(date.replace("-", ""))

    cxl_flag = False
    while True:
        v = requests.put(f"https://api.aditya.diwakar.io/cot/{cot_class}/{cot_id}",
                         json=record, headers={"Authorization": os.getenv("API_TOKEN")})
        if v.status_code == 404:
            new_prod = {
                "id": cot_id,
                "name": name
            }
            requests.post(f"https://api.aditya.diwakar.io/cot/{cot_class}",
                          json=new_prod, headers={"Authorization": os.getenv("API_TOKEN")})
            continue
        elif v.status_code == 409:
            # 409 means this date already exists in the system, cancel the entire script
            cxl_flag = True
            break
        else:
            break

    if cxl_flag:
        print("Cancelling early because this data is not new!")
        break

    print(f"Succesfully added {name}'s entry for {date} ({cot_class})")


end = time.time()

print(f"FINISHED INGESTION THAT TOOK {round(end-start, 2)} SECONDS")
