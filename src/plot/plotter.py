import json
import matplotlib.pyplot as plt
import numpy as np
import sys
import os


def plot_entropy(data_file):
    with open(data_file) as f:
        data = json.load(f)

    lengths = np.array(data['lengths'])
    entropies = np.array(data['entropies'])

    # Линейная регрессия
    coeffs = np.polyfit(lengths, entropies, 1)
    regression_line = np.poly1d(coeffs)

    plt.figure(figsize=(10, 6))

    # Столбчатая диаграмма
    bars = plt.bar(lengths, entropies, width=0.5,
                   alpha=0.7, label='Actual Entropy')

    x_fit = np.linspace(min(lengths)-0.5, max(lengths)+0.5, 100)
    plt.plot(x_fit, regression_line(x_fit), 'r--',
             label=f'Linear Regression (y = {coeffs[0]:.2f}x + {coeffs[1]:.2f})')

    # Настройки графика
    plt.title("Entropy vs Code Length with Regression", pad=20)
    plt.xlabel("Code Length (bytes)", labelpad=10)
    plt.ylabel("Entropy (bits)", labelpad=10)
    plt.xticks(lengths)
    plt.grid(axis='y', alpha=0.3)

    for bar in bars:
        height = bar.get_height()
        plt.text(bar.get_x() + bar.get_width()/2., height,
                 f'{height:.2f}',
                 ha='center', va='bottom')

    plt.legend(loc='upper left', frameon=True)
    plt.tight_layout()

    output_path = data['output']
    os.makedirs(os.path.dirname(output_path), exist_ok=True)
    plt.savefig(output_path, dpi=120)
    plt.close()


if __name__ == "__main__":
    plot_entropy(sys.argv[1])
