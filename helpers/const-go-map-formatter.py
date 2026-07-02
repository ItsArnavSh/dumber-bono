# Paste the text to be formatted into a golang map format and it will automatically do so
# Works as long as there are no numbers in the names


title = "CurrentSector"
text = """
0 = sector1, 1 = sector2, 2 = sector3
"""


pairs = {}
num_mode = True
numbuffer = ""
namebuffer = ""
for ch in text:
    if ch.isdecimal() or ch == "-":
        if len(namebuffer.strip()) != 0:
            pairs[numbuffer] = namebuffer.strip()
            namebuffer = ""
            numbuffer = ""
        numbuffer += ch
    else:
        if ch not in ["=", ","]:
            namebuffer += ch
pairs[numbuffer] = namebuffer.strip()
format_string = f"var {title} = map[int16]string{{"
for key, value in pairs.items():
    format_string += f'{key}: "{value}",\n'


format_string += "}"
print(format_string)
