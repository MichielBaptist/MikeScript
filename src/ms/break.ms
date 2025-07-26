
0 => x;
while x < 3 {
    0 => y;
    while y < 3 {
        "y", y >>= print;
        y + 1 -> y;

        if y > 1 {
            break;
        }
    }
    "x", x >>= print;
    x + 1 -> x;
}