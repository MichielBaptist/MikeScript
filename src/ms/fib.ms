function (int n) >> fib -> int {

    if (n <= 1) {
        return 1;
    }
    n - 1 >>= fib => fib_1;
    n - 2 >>= fib => fib_2;
    return fib_1 + fib_2;
}