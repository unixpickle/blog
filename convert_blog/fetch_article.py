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
        img_name = img_name.split("?")[0]

        img_path = None
        while img_path is None or os.path.exists(img_path):
            if img_path is not None:
                img_name = "1-" + img_name
            img_path = os.path.join(args.output_dir, "img", img_name)
        assert not os.path.exists(img_path)

        with open(img_path, "wb") as img_file:
            img_file.write(img_data)

        img_tag["src"] = f"img/{img_name}"
        for attr in ["decoding", "srcset"]:
            if attr in img_tag.attrs:
                del img_tag.attrs[attr]

        parent_a_tag = img_tag.find_parent("a")
        if parent_a_tag:
            parent_a_tag["href"] = img_tag["src"]

    print("cleaning up unused content...")
    for footer in soup.find_all("footer"):
        footer.decompose()

    html_content = f"<html><head></head><body>{article}</body></html>".replace(
        "https://blog.aqnichol.com/", "/"
    )

    with open(
        os.path.join(args.output_dir, "index.html"), "w", encoding="utf-8"
    ) as html_file:
        html_file.write(html_content)


if __name__ == "__main__":
    main()
