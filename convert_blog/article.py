import argparse
import os
from urllib.parse import urljoin

import requests
from bs4 import BeautifulSoup


def main():
    parser = argparse.ArgumentParser()
    parser.add_argument("--url", required=True, type=str)
    parser.add_argument("--output_dir", required=True, type=str)
    args = parser.parse_args()

    print("fetching page...")
    data = requests.get(args.url).text

    print("parsing...")
    soup = BeautifulSoup(data, "html.parser")

    article = soup.find("article")
    if not article:
        raise ValueError("No <article> tag found on the page.")

    article = soup.find("article")
    if not article:
        raise ValueError("No <article> tag found on the page.")

    os.makedirs(os.path.join(args.output_dir, "img"), exist_ok=True)

    print("downloading images...")
    for img_tag in article.find_all("img"):
        img_url = img_tag.get("src")
        if not img_url:
            continue

        img_url = urljoin(args.url, img_url)

        img_data = requests.get(img_url).content

        img_name = os.path.basename(img_url)
        img_path = os.path.join(args.output_dir, "img", img_name)
        with open(img_path, "wb") as img_file:
            img_file.write(img_data)

        img_tag["src"] = f"img/{img_name}"

        html_content = f"<html><body>{article}</body></html>"

    with open(
        os.path.join(args.output_dir, "index.html"), "w", encoding="utf-8"
    ) as html_file:
        html_file.write(html_content)


if __name__ == "__main__":
    main()
