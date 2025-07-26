function () >> f {
    0 => x;
    while x < 10 {
        if x == 5 {
            return;
        }
        "x", x >>= print;
        x + 1 -> x;
    }
    return x;
}

=f;