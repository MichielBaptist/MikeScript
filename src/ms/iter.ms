

[10]int{} => arr;

function (int x) >> f -> int{
    "call", x >>= print;
    return x + 1;
}

function (int x) >> even -> bool{
    return (x % 2) == 0;
}

arr .>>= even => res;
