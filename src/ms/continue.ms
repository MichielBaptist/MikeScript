var int x;

while x < 3 {
    "x:", x, "x*x:", x * x >>= print;

    if x % 2 == 0 {
        "x is even" >>= print;
    } else {
        "x is odd" >>= print;

        var int y;
        0 -> y;
        while y < 3 {
            "y:", y >>= print;
            if y > 1 {
                break;
            }
            y + 1 -> y;
        }
    }

    x + 1 -> x;
}