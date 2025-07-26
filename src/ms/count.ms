function (int x) >> count -> int {
    x >>= print;
    if x > 0 {
        x - 1 >>= count;
    }
    return x;
}

5 >>= count;