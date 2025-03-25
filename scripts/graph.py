import json
import sys
import matplotlib.pyplot as plt


SETTINGS = {
    "compression": {
        "xlabel": "Buffer Size (bytes)",
        "ylabel": "Compression Ratio",
        "title": "Dependence of Compression Ratio on Buffer Size",
        "output": "results/buffer_plot.png"
    },
    "entropy": {
        "xlabel": "Block Size (bytes)",
        "ylabel": "Average Entropy (bits per byte)",
        "title": "Dependence of Entropy on Block Size",
        "output": "results/entropy_plot.png"
    }
}


def load_data(file_path):
    with open(file_path, 'r') as f:
        return json.load(f)


def plot_data(buffer_sizes, values, xlabel, ylabel, title, output_file):
    plt.figure(figsize=(10, 6))
    plt.plot(buffer_sizes, values, marker='o', linestyle='-')
    plt.xlabel(xlabel)
    plt.ylabel(ylabel)
    plt.title(title)
    plt.grid(True, which='both', linestyle='--', linewidth=0.5)
    plt.savefig(output_file, dpi=300, bbox_inches='tight')


def get_data_type(file_name):
    if "entropy.json" in file_name:
        return "entropy"
    elif "buffer.json" in file_name:
        return "compression"


def main():
    file_name = sys.argv[1]

    data_type = get_data_type(file_name)

    settings = SETTINGS[data_type]
    xlabel = settings["xlabel"]
    ylabel = settings["ylabel"]
    title = settings["title"]
    output_file = settings["output"]

    results = load_data(file_name)

    buffer_sizes = [int(k) for k in results.keys()]
    values = [v for v in results.values()]

    sorted_indices = sorted(range(len(buffer_sizes)), key=lambda i: buffer_sizes[i])
    buffer_sizes = [buffer_sizes[i] for i in sorted_indices]
    values = [values[i] for i in sorted_indices]

    plot_data(buffer_sizes, values, xlabel, ylabel, title, output_file)


if __name__ == "__main__":
    main()