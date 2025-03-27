from PIL import Image
import numpy as np


def png_to_raw(image_path, output_path):
    image = Image.open(image_path)
    if image.mode in ('RGBA', 'LA') or (image.mode == 'P' and 'transparency' in image.info):
        image = image.convert('RGB')

    raw_pixels = np.array(image)
    raw_data = raw_pixels.tobytes()

    with open(output_path, 'wb') as f:
        f.write(raw_data)


if __name__ == '__main__':
    png_to_raw(
        "/home/kudrix/GolandProjects/AlgosLab1Sem4/datasets/image/gray.png",
        "/home/kudrix/GolandProjects/AlgosLab1Sem4/datasets/image/raw.gray"
    )

    png_to_raw(
        "/home/kudrix/GolandProjects/AlgosLab1Sem4/datasets/image/color.png",
        "/home/kudrix/GolandProjects/AlgosLab1Sem4/datasets/image/raw.color"
    )

    png_to_raw(
        "/home/kudrix/GolandProjects/AlgosLab1Sem4/datasets/image/bw.png",
        "/home/kudrix/GolandProjects/AlgosLab1Sem4/datasets/image/raw.bw"
    )

