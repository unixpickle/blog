import argparse
import os
from urllib.parse import urljoin

import requests
from bs4 import BeautifulSoup


def main():
    parser = argparse.ArgumentParser()
    parser.add_argument(
        "--latest_article",
        default="https://blog.aqnichol.com/2024/01/01/vpn-tunneling-to-share-streaming-services/",
    )
    args = parser.parse_args()

    url = args.latest_article
    while True:
        data = requests.get(url).text
        soup = BeautifulSoup(data, "html.parser")
        try:
            link = soup.find("div", class_="nav-previous").find("a")["href"]
        except:
            break
        url = link
        print(link)


if __name__ == "__main__":
    main()
