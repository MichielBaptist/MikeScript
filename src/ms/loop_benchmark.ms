
function (int u, int nl, int r) >> loops -> int {

    // alloc int array
    [nl]int{} => a;
    
    0 => i;
    while i < nl {
        0 => j;
        while j < nl {
            a[i] + j % u -> a[i];
            j + 1 -> j;
        }
        a[i] + r -> a[i];
        i + 1 -> i;
    }
    return a[0];
}

123, 10000, 234 >>= loops;
