from PIL import Image

def image_to_ascii(path, threshold=128, block="█", space=" "):
    img = Image.open(path).convert("L")  # convert to grayscale
    # Optionally resize for terminal
    # img = img.resize((img.width // 2, img.height // 2))
    pixels = img.load()
    w, h = img.size
    lines = []
    for y in range(h):
        line = ""
        for x in range(w):
            if pixels[x, y] < threshold:
                line += block
            else:
                line += space
        lines.append(line.rstrip())
    return "\n".join(lines)

# Example usage:
ascii_art = image_to_ascii("NexaAuto_Pixel.png", threshold=100)
print(ascii_art)
