int x;

while x < 10 {
    "x:", x, "x*x:", x * x >> print;

    if x % 2 == 0 {
        "x is even" >> print;
    } else {
        "x is odd" >> print;

        int y;
        0 -> y;
        while y < 10 {
            "y:", y >> print;
            if y > 5 {
                break;
            }
            y + 1 -> y;
        }
    }

    x + 1 -> x;
}
