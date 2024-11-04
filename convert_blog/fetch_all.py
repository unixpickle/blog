import argparse
import os
import subprocess
import sys

DIR_PATH = os.path.dirname(os.path.abspath(__file__))


def main():
    parser = argparse.ArgumentParser()
    parser.add_argument("--output_dir", default="articles")
    args = parser.parse_args()

    articles = [
        "https://blog.aqnichol.com/2024/01/01/vpn-tunneling-to-share-streaming-services/",
        "https://blog.aqnichol.com/2023/05/20/representing-3d-models-as-decision-trees/",
        "https://blog.aqnichol.com/2022/12/31/large-scale-vehicle-classification/",
        "https://blog.aqnichol.com/2021/04/12/data-and-machines/",
        "https://blog.aqnichol.com/2020/03/04/vq-draw-a-new-generative-model/",
        "https://blog.aqnichol.com/2020/01/18/research-didnt-pan-out/",
        "https://blog.aqnichol.com/2019/07/24/competing-in-the-obstacle-tower-challenge/",
        "https://blog.aqnichol.com/2019/04/03/prierarchy-implicit-hierarchies/",
        "https://blog.aqnichol.com/2018/12/24/solving-murder-with-go/",
        "https://blog.aqnichol.com/2017/12/23/what-i-dont-know/",
        "https://blog.aqnichol.com/2017/10/31/adversarial-traintest-splits/",
        "https://blog.aqnichol.com/2017/08/30/decision-trees-as-rl-policies/",
        "https://blog.aqnichol.com/2017/07/04/keeping-tabs-on-all-my-neural-networks/",
        "https://blog.aqnichol.com/2017/06/11/why-im-remaking-openai-universe/",
        "https://blog.aqnichol.com/2017/04/15/the-meta-learning-quest-part-1/",
        "https://blog.aqnichol.com/2017/04/05/slice-aliasing/",
        "https://blog.aqnichol.com/2017/03/08/the-bug-that-wasted-a-month-of-gpu-time/",
        "https://blog.aqnichol.com/2017/03/01/random-fun-with-linear-svms/",
        "https://blog.aqnichol.com/2017/02/14/can-neural-networks-learn-to-spell/",
        "https://blog.aqnichol.com/2017/02/12/ancient-philosophy-as-a-classification-problem/",
    ]
    for article in articles:
        subprocess.check_call(
            [
                sys.executable,
                os.path.join(DIR_PATH, "fetch_article.py"),
                "--url",
                article,
                "--output_dir",
                os.path.join(
                    args.output_dir, article.replace("https://blog.aqnichol.com/", "")
                ),
            ]
        )


if __name__ == "__main__":
    main()
